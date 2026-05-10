package ssh

import (
	"context"
	"net"
	"testing"
	"time"

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
