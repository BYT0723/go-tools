package sshx

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"net"
	"testing"
	"time"

	"github.com/BYT0723/go-tools/transport/connmux"
	"github.com/stretchr/testify/assert"
	xssh "golang.org/x/crypto/ssh"
)

type mockConnMeta struct {
	user string
}

func (m *mockConnMeta) User() string          { return m.user }
func (m *mockConnMeta) SessionID() []byte     { return nil }
func (m *mockConnMeta) ClientVersion() []byte { return nil }
func (m *mockConnMeta) ServerVersion() []byte { return nil }
func (m *mockConnMeta) RemoteAddr() net.Addr  { return nil }
func (m *mockConnMeta) LocalAddr() net.Addr   { return nil }

func TestNewServer(t *testing.T) {
	s := NewServer()
	assert.NotNil(t, s)
	assert.Equal(t, ":2222", s.Addr())
}

func TestWithAddr(t *testing.T) {
	s := NewServer(WithAddr(":9999"))
	assert.Equal(t, ":9999", s.addr)
}

func TestWithUser(t *testing.T) {
	s := NewServer(
		WithUser("root", "secret"),
		WithAddr(":0"),
	)
	assert.NotNil(t, s.config.PasswordCallback)

	perms, err := s.config.PasswordCallback(
		&mockConnMeta{user: "root"},
		[]byte("secret"),
	)
	assert.Nil(t, err)
	assert.Nil(t, perms)
}

func TestWithPasswordAuth(t *testing.T) {
	ch := make(chan string, 1)
	s := NewServer(
		WithAddr(":0"),
		WithPasswordAuth(func(user, password string) bool {
			ch <- user + ":" + password
			return true
		}),
	)

	perms, err := s.config.PasswordCallback(
		&mockConnMeta{user: "admin"},
		[]byte("pass"),
	)
	assert.Nil(t, err)
	assert.Nil(t, perms)
	assert.Equal(t, "admin:pass", <-ch)
}

func TestStartStop(t *testing.T) {
	s := NewServer(WithAddr(":0"))

	ctx, cf := context.WithTimeout(context.Background(), 5*time.Second)
	defer cf()

	err := s.Start(ctx)
	assert.Nil(t, err)

	time.Sleep(50 * time.Millisecond)
	err = s.Stop(ctx)
	assert.Nil(t, err)
}

func TestStartTwice(t *testing.T) {
	s := NewServer(WithAddr(":0"))
	ctx := context.Background()

	err := s.Start(ctx)
	assert.Nil(t, err)

	err = s.Start(ctx)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "already running")

	s.Stop(ctx)
}

func TestSSHConnection(t *testing.T) {
	hostKey, err := GenerateHostKey()
	assert.Nil(t, err)

	ctx, cf := context.WithTimeout(context.Background(), 10*time.Second)
	defer cf()

	s := NewServer(
		WithAddr("127.0.0.1:0"),
		WithHostKey(hostKey),
		WithPasswordAuth(func(user, password string) bool {
			return user == "test" && password == "pass"
		}),
		WithHandler(func(sessConn *xssh.ServerConn, ch xssh.Channel, reqs <-chan *xssh.Request) {
			defer ch.Close()
			for req := range reqs {
				switch req.Type {
				case "shell":
					req.Reply(true, nil)
					ch.Write([]byte("hello\n"))
				case "exec":
					req.Reply(true, nil)
					ch.Write([]byte("ok\n"))
					ch.SendRequest("exit-status", false, xssh.Marshal(&struct{ Status uint32 }{0}))
					return
				default:
					req.Reply(false, nil)
				}
			}
		}),
	)

	err = s.Start(ctx)
	assert.Nil(t, err)
	defer s.Stop(ctx)

	time.Sleep(100 * time.Millisecond)
	addr := s.Addr()

	config := &xssh.ClientConfig{
		User: "test",
		Auth: []xssh.AuthMethod{
			xssh.Password("pass"),
		},
		HostKeyCallback: xssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	client, err := xssh.Dial("tcp", addr, config)
	assert.Nil(t, err)
	defer client.Close()

	sess, err := client.NewSession()
	assert.Nil(t, err)
	defer sess.Close()

	output, err := sess.Output("test")
	assert.Nil(t, err)
	assert.Equal(t, "ok\n", string(output))
}

func TestWithPublicKeyAuth(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(t, err)
	signer, err := xssh.NewSignerFromKey(key)
	assert.Nil(t, err)

	ch := make(chan bool, 1)
	s := NewServer(
		WithAddr(":0"),
		WithPublicKeyAuth(func(conn xssh.ConnMetadata, key xssh.PublicKey) bool {
			ch <- true
			return true
		}),
	)

	perms, err := s.config.PublicKeyCallback(
		&mockConnMeta{user: "test"},
		signer.PublicKey(),
	)
	assert.Nil(t, err)
	assert.Nil(t, perms)
	assert.True(t, <-ch)
}

func TestWithPublicKeyAuthReject(t *testing.T) {
	s := NewServer(
		WithAddr(":0"),
		WithPublicKeyAuth(func(conn xssh.ConnMetadata, key xssh.PublicKey) bool {
			return false
		}),
	)

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(t, err)
	signer, err := xssh.NewSignerFromKey(key)
	assert.Nil(t, err)

	_, err = s.config.PublicKeyCallback(
		&mockConnMeta{user: "test"},
		signer.PublicKey(),
	)
	assert.NotNil(t, err)
}

func TestWithShellPath(t *testing.T) {
	s := NewServer(
		WithShellPath("/bin/zsh"),
	)
	assert.Equal(t, "/bin/zsh", s.shellPath)
}

func TestSrvxLifecycle(t *testing.T) {
	s := NewServer(WithAddr(":0"))
	assert.Equal(t, "sshx", s.Name())

	ctx, cf := context.WithTimeout(context.Background(), 5*time.Second)
	defer cf()

	err := s.Init(ctx)
	assert.Nil(t, err)

	// Run blocks; call in goroutine and stop via context
	runCtx, runCf := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- s.Run(runCtx)
	}()
	time.Sleep(100 * time.Millisecond)
	runCf()

	err = <-done
	assert.Nil(t, err)

	err = s.Destroy(ctx)
	assert.Nil(t, err)
}

func TestExecCommand(t *testing.T) {
	hostKey, err := GenerateHostKey()
	assert.Nil(t, err)

	ctx, cf := context.WithTimeout(context.Background(), 10*time.Second)
	defer cf()

	s := NewServer(
		WithAddr("127.0.0.1:0"),
		WithHostKey(hostKey),
		WithUser("test", "pass"),
	)

	err = s.Start(ctx)
	assert.Nil(t, err)
	defer s.Stop(ctx)
	time.Sleep(100 * time.Millisecond)

	config := &xssh.ClientConfig{
		User: "test",
		Auth: []xssh.AuthMethod{
			xssh.Password("pass"),
		},
		HostKeyCallback: xssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	client, err := xssh.Dial("tcp", s.Addr(), config)
	assert.Nil(t, err)
	defer client.Close()

	sess, err := client.NewSession()
	assert.Nil(t, err)
	defer sess.Close()

	output, err := sess.Output("echo hello")
	assert.Nil(t, err)
	assert.Contains(t, string(output), "hello")
}

func TestExecCommandExitCode(t *testing.T) {
	hostKey, err := GenerateHostKey()
	assert.Nil(t, err)

	ctx, cf := context.WithTimeout(context.Background(), 10*time.Second)
	defer cf()

	s := NewServer(
		WithAddr("127.0.0.1:0"),
		WithHostKey(hostKey),
		WithUser("test", "pass"),
	)

	err = s.Start(ctx)
	assert.Nil(t, err)
	defer s.Stop(ctx)
	time.Sleep(100 * time.Millisecond)

	config := &xssh.ClientConfig{
		User: "test",
		Auth: []xssh.AuthMethod{
			xssh.Password("pass"),
		},
		HostKeyCallback: xssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	client, err := xssh.Dial("tcp", s.Addr(), config)
	assert.Nil(t, err)
	defer client.Close()

	sess, err := client.NewSession()
	assert.Nil(t, err)
	defer sess.Close()

	err = sess.Run("exit 42")
	assert.NotNil(t, err)

	if exitErr, ok := err.(*xssh.ExitError); ok {
		assert.Equal(t, 42, exitErr.ExitStatus())
	} else {
		t.Fatalf("expected ExitError, got %v", err)
	}
}

func TestDirectTcpipForwarding(t *testing.T) {
	// Start a target TCP server
	targetLn, err := net.Listen("tcp", "127.0.0.1:0")
	assert.Nil(t, err)
	defer targetLn.Close()

	go func() {
		conn, err := targetLn.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		conn.Write([]byte("from-target"))
	}()

	hostKey, err := GenerateHostKey()
	assert.Nil(t, err)

	ctx, cf := context.WithTimeout(context.Background(), 10*time.Second)
	defer cf()

	s := NewServer(
		WithAddr("127.0.0.1:0"),
		WithHostKey(hostKey),
		WithUser("test", "pass"),
	)

	err = s.Start(ctx)
	assert.Nil(t, err)
	defer s.Stop(ctx)
	time.Sleep(100 * time.Millisecond)

	config := &xssh.ClientConfig{
		User: "test",
		Auth: []xssh.AuthMethod{
			xssh.Password("pass"),
		},
		HostKeyCallback: xssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	client, err := xssh.Dial("tcp", s.Addr(), config)
	assert.Nil(t, err)
	defer client.Close()

	conn, err := client.Dial("tcp", targetLn.Addr().String())
	assert.Nil(t, err)
	defer conn.Close()

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	assert.Nil(t, err)
	assert.Equal(t, "from-target", string(buf[:n]))
}

func TestTcpipForwardRequest(t *testing.T) {
	hostKey, err := GenerateHostKey()
	assert.Nil(t, err)

	ctx, cf := context.WithTimeout(context.Background(), 10*time.Second)
	defer cf()

	s := NewServer(
		WithAddr("127.0.0.1:0"),
		WithHostKey(hostKey),
		WithUser("test", "pass"),
	)

	err = s.Start(ctx)
	assert.Nil(t, err)
	defer s.Stop(ctx)
	time.Sleep(100 * time.Millisecond)

	config := &xssh.ClientConfig{
		User: "test",
		Auth: []xssh.AuthMethod{
			xssh.Password("pass"),
		},
		HostKeyCallback: xssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	client, err := xssh.Dial("tcp", s.Addr(), config)
	assert.Nil(t, err)
	defer client.Close()

	// Request remote forwarding
	listener, err := client.Listen("tcp", "127.0.0.1:0")
	assert.Nil(t, err)
	defer listener.Close()

	resultCh := make(chan string, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		buf := make([]byte, 4096)
		n, _ := io.ReadFull(conn, buf[:11])
		resultCh <- string(buf[:n])
	}()

	// Connect to the forwarded port
	time.Sleep(50 * time.Millisecond)
	conn, err := net.Dial("tcp", listener.Addr().String())
	assert.Nil(t, err)
	conn.Write([]byte("reverse-fwd"))
	conn.Close()

	select {
	case result := <-resultCh:
		assert.Equal(t, "reverse-fwd", result)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for forwarded connection")
	}
}

func TestShellSession(t *testing.T) {
	hostKey, err := GenerateHostKey()
	assert.Nil(t, err)

	ctx, cf := context.WithTimeout(context.Background(), 15*time.Second)
	defer cf()

	s := NewServer(
		WithAddr("127.0.0.1:0"),
		WithHostKey(hostKey),
		WithUser("test", "pass"),
	)

	err = s.Start(ctx)
	assert.Nil(t, err)
	defer s.Stop(ctx)
	time.Sleep(100 * time.Millisecond)

	config := &xssh.ClientConfig{
		User: "test",
		Auth: []xssh.AuthMethod{
			xssh.Password("pass"),
		},
		HostKeyCallback: xssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	client, err := xssh.Dial("tcp", s.Addr(), config)
	assert.Nil(t, err)
	defer client.Close()

	sess, err := client.NewSession()
	assert.Nil(t, err)
	defer sess.Close()

	err = sess.RequestPty("xterm", 40, 80, xssh.TerminalModes{})
	assert.Nil(t, err)

	stdin, err := sess.StdinPipe()
	assert.Nil(t, err)
	stdout, err := sess.StdoutPipe()
	assert.Nil(t, err)

	err = sess.Shell()
	assert.Nil(t, err)

	stdin.Write([]byte("echo SHELL_OK\n"))
	stdin.Write([]byte("exit\n"))

	buf := make([]byte, 4096)
	n, err := stdout.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("read error: %v", err)
	}
	output := string(buf[:n])
	assert.Contains(t, output, "SHELL_OK")
}

func TestMuxSSHIntegration(t *testing.T) {
	key, err := GenerateHostKey()
	assert.NoError(t, err)

	mux := connmux.NewMux()

	sshSrv := NewServer(
		WithHostKey(key),
		WithPasswordAuth(func(user, password string) bool {
			return user == "test" && password == "test"
		}),
	)

	mux.Route("ssh", sshSrv)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- mux.Run(ctx)
	}()

	time.Sleep(150 * time.Millisecond)
	addr := mux.Addr()

	clientCfg := &xssh.ClientConfig{
		User:            "test",
		Auth:            []xssh.AuthMethod{xssh.Password("test")},
		HostKeyCallback: xssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	assert.NoError(t, err)

	sshConn, chans, reqs, err := xssh.NewClientConn(conn, addr, clientCfg)
	assert.NoError(t, err)

	client := xssh.NewClient(sshConn, chans, reqs)

	session, err := client.NewSession()
	assert.NoError(t, err)

	output, err := session.Output("echo CONNMUX_OK")
	assert.NoError(t, err)
	assert.Contains(t, string(output), "CONNMUX_OK")

	session.Close()
	client.Close()
	sshConn.Close()
	conn.Close()

	cancel()
	select {
	case err = <-errCh:
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for mux shutdown")
	}
}
