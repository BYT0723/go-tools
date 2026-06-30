package frpx

import (
	"bytes"
	"context"
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	c := NewClient()
	assert.NotNil(t, c)
	assert.Equal(t, "frpx", c.Name())
	assert.Equal(t, "127.0.0.1:7000", c.addr)
}

func TestWithServerAddr(t *testing.T) {
	c := NewClient(WithServerAddr("frp.example.com:7000"))
	assert.Equal(t, "frp.example.com:7000", c.addr)
}

func TestWithToken(t *testing.T) {
	c := NewClient(WithToken("my-token"))
	assert.Equal(t, "my-token", c.token)
}

func TestAuthKey(t *testing.T) {
	k := authKey("test-token", 1234567890)
	assert.NotEmpty(t, k)
	assert.Len(t, k, 32)
}

func TestProxyConfigToNewProxy(t *testing.T) {
	cfg := ProxyConfig{
		Name:          "web",
		Type:          "http",
		LocalAddr:     "127.0.0.1:8080",
		CustomDomains: []string{"app.example.com"},
	}
	msg := cfg.toNewProxy("testuser")
	assert.Equal(t, "testuser.web", msg.ProxyName)
	assert.Equal(t, "http", msg.ProxyType)
	assert.Equal(t, []string{"app.example.com"}, msg.CustomDomains)
}

func TestClientLifecycle(t *testing.T) {
	mock := newMockFrps()
	addr, err := mock.listen()
	assert.Nil(t, err)
	defer mock.close()

	c := NewClient(
		WithServerAddr(addr),
		WithToken("test-token"),
		WithUser("test"),
		WithHeartbeatInterval(500*time.Millisecond),
		WithReconnectBackoff(100*time.Millisecond),
		WithProxy(ProxyConfig{
			Name:       "tcp-proxy",
			Type:       "tcp",
			LocalAddr:  "127.0.0.1:9999",
			RemotePort: 8080,
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = c.Init(ctx)
	assert.Nil(t, err)

	done := make(chan error, 1)
	go func() {
		done <- c.Run(ctx)
	}()

	time.Sleep(500 * time.Millisecond)

	c.mu.Lock()
	assert.NotEmpty(t, c.runID)
	c.mu.Unlock()

	cancel()
	err = <-done
	assert.Nil(t, err)

	err = c.Destroy(ctx)
	assert.Nil(t, err)
}

func TestClientReconnect(t *testing.T) {
	mock := newMockFrps()
	addr, err := mock.listen()
	assert.Nil(t, err)

	c := NewClient(
		WithServerAddr(addr),
		WithToken("test-token"),
		WithReconnectBackoff(50*time.Millisecond),
		WithHeartbeatInterval(100*time.Millisecond),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go c.Run(ctx)

	time.Sleep(300 * time.Millisecond)
	c.mu.Lock()
	assert.NotEmpty(t, c.runID)
	c.mu.Unlock()

	mock.close()
	time.Sleep(500 * time.Millisecond)

	cancel()
	time.Sleep(100 * time.Millisecond)
}

func TestV2FrameCodec(t *testing.T) {
	var buf bytes.Buffer
	err := writeV2Frame(&buf, FrameTypeClientHello, []byte(`{"hello":"world"}`))
	assert.Nil(t, err)

	typ, body, err := readV2Frame(&buf)
	assert.Nil(t, err)
	assert.Equal(t, FrameTypeClientHello, typ)
	assert.Equal(t, `{"hello":"world"}`, string(body))
}

func TestV2MessageCodec(t *testing.T) {
	var buf bytes.Buffer
	err := writeV2Msg(&buf, msgTypePing, Ping{Timestamp: 123})
	assert.Nil(t, err)

	typ, body, err := readV2Msg(&buf)
	assert.Nil(t, err)
	assert.Equal(t, msgTypePing, typ)
	assert.Contains(t, string(body), "123")
}

func TestAEADStream(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	key1 := make([]byte, 32)
	key2 := make([]byte, 32)
	for i := range key1 {
		key1[i] = byte(i)
		key2[i] = byte(i + 32)
	}

	done := make(chan []byte, 1)
	go func() {
		s := newAEADStream(server, key1, key2)
		buf := make([]byte, 11)
		io.ReadFull(s, buf)
		done <- buf
	}()

	c := newAEADStream(client, key2, key1)
	c.Write([]byte("hello world"))

	result := <-done
	assert.Equal(t, "hello world", string(result))
}
