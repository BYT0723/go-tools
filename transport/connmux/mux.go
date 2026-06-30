package connmux

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type ListenedService interface {
	Name() string
	Init(ctx context.Context) error
	Run(ctx context.Context) error
	Destroy(ctx context.Context) error
	SetListener(net.Listener)
	Match() Matcher
}

type Mux struct {
	mu            sync.Mutex
	addr          string
	sniffSize     int
	sniffDeadline time.Duration
	listener      net.Listener
	running       bool
	routes        []*routeEntry
}

type routeEntry struct {
	name    string
	matcher Matcher
	vl      *VirtualListener
	svc     ListenedService
}

func NewMux(opts ...Option) *Mux {
	m := &Mux{
		addr:          ":0",
		sniffSize:     256,
		sniffDeadline: 5 * time.Second,
	}
	for _, o := range opts {
		o(m)
	}
	return m
}

func (m *Mux) Route(name string, svc ListenedService) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.running {
		panic("connmux: Route called after Run")
	}
	m.routes = append(m.routes, &routeEntry{
		name:    name,
		matcher: svc.Match(),
		svc:     svc,
	})
}

func (m *Mux) Name() string {
	return "connmux"
}

func (m *Mux) Addr() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.listener != nil {
		return m.listener.Addr().String()
	}
	return m.addr
}

func (m *Mux) Init(_ context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.sniffSize <= 0 || m.sniffSize > 65536 {
		return fmt.Errorf("connmux: sniffSize must be 1-65536, got %d", m.sniffSize)
	}
	if len(m.routes) == 0 {
		return fmt.Errorf("connmux: at least one route required")
	}
	return nil
}

func (m *Mux) Run(ctx context.Context) error {
	return m.runOnce(ctx)
}

func (m *Mux) Destroy(_ context.Context) error {
	m.mu.Lock()
	m.running = false
	if m.listener != nil {
		m.listener.Close()
	}
	m.mu.Unlock()
	return nil
}

func (m *Mux) runOnce(ctx context.Context) error {
	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return fmt.Errorf("connmux: already running")
	}
	m.running = true

	if len(m.routes) == 0 {
		m.running = false
		m.mu.Unlock()
		return fmt.Errorf("connmux: no routes registered")
	}

	deriveCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	lc := net.ListenConfig{}
	l, err := lc.Listen(ctx, "tcp", m.addr)
	if err != nil {
		m.running = false
		m.mu.Unlock()
		return err
	}
	m.listener = l
	addr := l.Addr()

	for i, route := range m.routes {
		route.vl = newVirtualListener(addr)
		route.svc.SetListener(route.vl)
		if err := route.svc.Init(deriveCtx); err != nil {
			m.cleanupRoutes(i, deriveCtx)
			m.running = false
			m.mu.Unlock()
			l.Close()
			return fmt.Errorf("connmux: init %s: %w", route.name, err)
		}
	}

	m.mu.Unlock()

	var svcWg sync.WaitGroup
	svcDone := make(chan struct{})

	for _, route := range m.routes {
		svcWg.Add(1)
		go func(r *routeEntry) {
			defer svcWg.Done()
			_ = r.svc.Run(deriveCtx)
		}(route)
	}
	go func() {
		svcWg.Wait()
		close(svcDone)
	}()

	var serveWg sync.WaitGroup
	serveWg.Add(1)
	go func() {
		defer serveWg.Done()
		m.serve(deriveCtx, l)
	}()

	<-ctx.Done()

	l.Close()
	serveWg.Wait()

	m.mu.Lock()
	for _, route := range m.routes {
		route.vl.Close()
	}
	m.mu.Unlock()

	cancel()
	<-svcDone

	m.shutdownRoutes(ctx)

	m.mu.Lock()
	m.running = false
	m.listener = nil
	m.mu.Unlock()

	return nil
}

func (m *Mux) serve(_ context.Context, l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}

		conn.SetReadDeadline(time.Now().Add(m.sniffDeadline))
		sniffed := make([]byte, m.sniffSize)
		n, readErr := io.ReadAtLeast(conn, sniffed, 1)
		if readErr != nil && n == 0 {
			conn.Close()
			continue
		}
		sniffed = sniffed[:n]

		var zeroTime time.Time
		conn.SetReadDeadline(zeroTime)

		wrapped := &replayConn{
			Conn:   conn,
			reader: io.MultiReader(bytes.NewReader(sniffed), conn),
		}

		matched := false
		for _, route := range m.routes {
			if route.matcher.Match(sniffed) {
				route.vl.push(wrapped)
				matched = true
				break
			}
		}

		if !matched {
			m.routes[len(m.routes)-1].vl.push(wrapped)
		}
	}
}

func (m *Mux) cleanupRoutes(fromIdx int, ctx context.Context) {
	for i := fromIdx - 1; i >= 0; i-- {
		_ = m.routes[i].svc.Destroy(ctx)
	}
}

func (m *Mux) shutdownRoutes(ctx context.Context) {
	m.mu.Lock()
	routes := make([]*routeEntry, len(m.routes))
	copy(routes, m.routes)
	m.mu.Unlock()

	for i := len(routes) - 1; i >= 0; i-- {
		_ = routes[i].svc.Destroy(ctx)
	}
}

type replayConn struct {
	net.Conn
	reader io.Reader
}

func (c *replayConn) Read(b []byte) (int, error) {
	return c.reader.Read(b)
}
