package connmux

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockService struct {
	mu       sync.Mutex
	name     string
	matcher  Matcher
	listener net.Listener
	inited   bool
	running  bool
	conns    []net.Conn
	initErr  error
	onRun    func(ctx context.Context) error
}

func (m *mockService) Name() string               { return m.name }
func (m *mockService) SetListener(l net.Listener) { m.listener = l }
func (m *mockService) Match() Matcher             { return m.matcher }
func (m *mockService) Init(ctx context.Context) error {
	if m.initErr == nil {
		m.mu.Lock()
		m.inited = true
		m.mu.Unlock()
	}
	return m.initErr
}
func (m *mockService) Destroy(ctx context.Context) error {
	m.mu.Lock()
	m.inited = false
	m.mu.Unlock()
	return nil
}
func (m *mockService) isInited() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.inited
}
func (m *mockService) Run(ctx context.Context) error {
	m.mu.Lock()
	m.running = true
	m.mu.Unlock()

	if m.onRun != nil {
		return m.onRun(ctx)
	}

	<-ctx.Done()
	m.mu.Lock()
	m.running = false
	m.mu.Unlock()
	return nil
}

func (m *mockService) isRunning() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.running
}

func newPrefix(p string) *bytePrefixMatcher {
	return &bytePrefixMatcher{prefix: []byte(p)}
}

func TestMuxRoute(t *testing.T) {
	mux := NewMux()
	svc := &mockService{name: "test", matcher: MatchSSH}

	mux.Route("test", svc)

	assert.Len(t, mux.routes, 1)
	assert.Equal(t, "test", mux.routes[0].name)
}

func TestMuxRoutePanicAfterRun(t *testing.T) {
	mux := NewMux()
	svc := &mockService{name: "test", matcher: MatchSSH}
	mux.Route("test", svc)

	mux.mu.Lock()
	mux.running = true
	mux.mu.Unlock()

	assert.Panics(t, func() {
		mux.Route("another", &mockService{name: "another"})
	})
}

func TestMuxInitNoRoutes(t *testing.T) {
	mux := NewMux()
	err := mux.Init(context.Background())
	assert.ErrorContains(t, err, "at least one route")
}

func TestMuxInitInvalidSniffSize(t *testing.T) {
	mux := NewMux(WithSniffSize(0))
	svc := &mockService{name: "test", matcher: MatchDefault}
	mux.Route("test", svc)

	err := mux.Init(context.Background())
	assert.ErrorContains(t, err, "sniffSize")
}

func TestMuxPartialInitRollback(t *testing.T) {
	mux := NewMux()
	svc1 := &mockService{name: "ok", matcher: MatchSSH}
	svc2 := &mockService{name: "fail", matcher: MatchDefault, initErr: fmt.Errorf("boom")}
	mux.Route("ok", svc1)
	mux.Route("fail", svc2)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := mux.Run(ctx)
	assert.ErrorContains(t, err, "boom")
	assert.False(t, svc1.isInited(), "svc1 should be cleaned up after rollback")
	assert.False(t, svc2.isInited())
}

func TestMuxLifecycle(t *testing.T) {
	mux := NewMux()
	svc := &mockService{name: "test", matcher: MatchDefault}
	mux.Route("test", svc)

	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		errCh <- mux.Run(ctx)
	}()

	time.Sleep(100 * time.Millisecond)

	assert.True(t, svc.isRunning(), "service should be running")

	cancel()

	select {
	case err := <-errCh:
		assert.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for mux to stop")
	}

	assert.False(t, svc.isRunning())
	assert.False(t, svc.isInited())
}

func TestMuxDispatch(t *testing.T) {
	mux := NewMux()

	var svc1Conn, svc2Conn net.Conn
	var svc1Mu, svc2Mu sync.Mutex

	svc1 := &mockService{name: "svc1", matcher: MatchSSH}
	svc2 := &mockService{name: "svc2", matcher: MatchHTTP1}

	svc1.onRun = func(ctx context.Context) error {
		go func() {
			for {
				conn, err := svc1.listener.Accept()
				if err != nil {
					return
				}
				svc1Mu.Lock()
				svc1Conn = conn
				svc1Mu.Unlock()
				conn.Close()
			}
		}()
		<-ctx.Done()
		return nil
	}

	svc2.onRun = func(ctx context.Context) error {
		go func() {
			for {
				conn, err := svc2.listener.Accept()
				if err != nil {
					return
				}
				svc2Mu.Lock()
				svc2Conn = conn
				svc2Mu.Unlock()
				conn.Close()
			}
		}()
		<-ctx.Done()
		return nil
	}

	mux.Route("svc1", svc1)
	mux.Route("svc2", svc2)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- mux.Run(ctx)
	}()

	time.Sleep(100 * time.Millisecond)

	addr := mux.Addr()

	sshConn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	sshConn.Write([]byte("SSH-2.0-test"))
	time.Sleep(50 * time.Millisecond)
	sshConn.Close()

	svc1Mu.Lock()
	assert.NotNil(t, svc1Conn, "SSH connection should be dispatched to svc1")
	svc1Mu.Unlock()

	httpConn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	httpConn.Write([]byte("GET / HTTP/1.1"))
	time.Sleep(50 * time.Millisecond)
	httpConn.Close()

	svc2Mu.Lock()
	assert.NotNil(t, svc2Conn, "HTTP connection should be dispatched to svc2")
	svc2Mu.Unlock()

	cancel()
	<-errCh
}

func TestMuxByteReplay(t *testing.T) {
	mux := NewMux()

	var (
		received    []byte
		receivedMu  sync.Mutex
	)
	svc := &mockService{name: "svc", matcher: MatchSSH}
	svc.onRun = func(ctx context.Context) error {
		go func() {
			for {
				conn, err := svc.listener.Accept()
				if err != nil {
					return
				}
				buf := make([]byte, 1024)
				n, _ := conn.Read(buf)
				receivedMu.Lock()
				received = buf[:n]
				receivedMu.Unlock()
				conn.Close()
			}
		}()
		<-ctx.Done()
		return nil
	}

	mux.Route("svc", svc)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go mux.Run(ctx)
	time.Sleep(100 * time.Millisecond)

	addr := mux.Addr()
	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	conn.Write([]byte("SSH-2.0-HELLO"))
	time.Sleep(50 * time.Millisecond)
	conn.Close()

	receivedMu.Lock()
	result := string(received)
	receivedMu.Unlock()
	assert.Equal(t, "SSH-2.0-HELLO", result, "sniffed bytes should be replayed to service")

	cancel()
}

func TestMuxConcurrent(t *testing.T) {
	mux := NewMux()

	var mu sync.Mutex
	var conns []net.Conn

	svc := &mockService{name: "svc", matcher: MatchDefault}
	svc.onRun = func(ctx context.Context) error {
		go func() {
			for {
				conn, err := svc.listener.Accept()
				if err != nil {
					return
				}
				mu.Lock()
				conns = append(conns, conn)
				mu.Unlock()
				conn.Close()
			}
		}()
		<-ctx.Done()
		return nil
	}

	mux.Route("svc", svc)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go mux.Run(ctx)
	time.Sleep(100 * time.Millisecond)

	addr := mux.Addr()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				return
			}
			conn.Write([]byte("SSH-2.0-test"))
			time.Sleep(10 * time.Millisecond)
			conn.Close()
		}()
	}
	wg.Wait()
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	count := len(conns)
	mu.Unlock()
	assert.GreaterOrEqual(t, count, 8, "at least 8/10 connections should be dispatched")

	cancel()
}

func TestMuxMultiHTTP(t *testing.T) {
	mux := NewMux()

	var svc1Conn, svc2Conn net.Conn
	var svc1Mu, svc2Mu sync.Mutex

	svc1 := &mockService{name: "api", matcher: newPrefix("GET /api")}
	svc1.onRun = func(ctx context.Context) error {
		go func() {
			for {
				conn, err := svc1.listener.Accept()
				if err != nil {
					return
				}
				buf := make([]byte, 1024)
				n, _ := conn.Read(buf)
				body := string(buf[:n])
				conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\napi:" + body))
				svc1Mu.Lock()
				svc1Conn = conn
				svc1Mu.Unlock()
				conn.Close()
			}
		}()
		<-ctx.Done()
		return nil
	}

	svc2 := &mockService{name: "web", matcher: MatchDefault}
	svc2.onRun = func(ctx context.Context) error {
		go func() {
			for {
				conn, err := svc2.listener.Accept()
				if err != nil {
					return
				}
				buf := make([]byte, 1024)
				n, _ := conn.Read(buf)
				body := string(buf[:n])
				conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nweb:" + body))
				svc2Mu.Lock()
				svc2Conn = conn
				svc2Mu.Unlock()
				conn.Close()
			}
		}()
		<-ctx.Done()
		return nil
	}

	mux.Route("api", svc1)
	mux.Route("web", svc2)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go mux.Run(ctx)
	time.Sleep(100 * time.Millisecond)

	addr := mux.Addr()

	apiConn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	apiConn.Write([]byte("GET /api/v1/users HTTP/1.1\r\n\r\n"))
	buf := make([]byte, 4096)
	n, _ := apiConn.Read(buf)
	apiResp := string(buf[:n])
	apiConn.Close()

	webConn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	webConn.Write([]byte("GET /index.html HTTP/1.1\r\n\r\n"))
	n, _ = webConn.Read(buf)
	webResp := string(buf[:n])
	webConn.Close()

	svc1Mu.Lock()
	c1 := svc1Conn
	svc1Mu.Unlock()
	assert.NotNil(t, c1, "api route should get connection")

	svc2Mu.Lock()
	c2 := svc2Conn
	svc2Mu.Unlock()
	assert.NotNil(t, c2, "web (default) should get connection")
	assert.Contains(t, apiResp, "api:")
	assert.Contains(t, webResp, "web:")

	cancel()
}

