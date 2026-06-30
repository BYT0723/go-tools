package connmux

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testAddr struct{}

func (testAddr) Network() string { return "tcp" }
func (testAddr) String() string  { return "virtual:0" }

type testConn struct{}

func (testConn) Read([]byte) (int, error)  { return 0, nil }
func (testConn) Write([]byte) (int, error) { return 0, nil }
func (testConn) Close() error              { return nil }
func (testConn) LocalAddr() net.Addr       { return testAddr{} }
func (testConn) RemoteAddr() net.Addr      { return testAddr{} }
func (testConn) SetDeadline(time.Time) error      { return nil }
func (testConn) SetReadDeadline(time.Time) error  { return nil }
func (testConn) SetWriteDeadline(time.Time) error { return nil }

func TestVirtualListenerAcceptPush(t *testing.T) {
	vl := newVirtualListener(testAddr{})

	go func() {
		vl.push(&testConn{})
	}()

	conn, err := vl.Accept()
	assert.NoError(t, err)
	assert.NotNil(t, conn)
}

func TestVirtualListenerClose(t *testing.T) {
	vl := newVirtualListener(testAddr{})

	go func() {
		time.Sleep(10 * time.Millisecond)
		vl.Close()
	}()

	_, err := vl.Accept()
	assert.Equal(t, ErrListenerClosed, err)
}

func TestVirtualListenerCloseIdempotent(t *testing.T) {
	vl := newVirtualListener(testAddr{})

	assert.NoError(t, vl.Close())
	assert.NoError(t, vl.Close())
	assert.NoError(t, vl.Close())
}

func TestVirtualListenerAddr(t *testing.T) {
	vl := newVirtualListener(testAddr{})
	assert.Equal(t, "virtual:0", vl.Addr().String())
}

func TestVirtualListenerPushClosed(t *testing.T) {
	vl := newVirtualListener(testAddr{})
	vl.Close()

	vl.push(&testConn{})
}

func TestVirtualListenerPushBuffered(t *testing.T) {
	vl := newVirtualListener(testAddr{})

	for i := 0; i < 5; i++ {
		vl.push(&testConn{})
	}

	for i := 0; i < 5; i++ {
		conn, err := vl.Accept()
		assert.NoError(t, err)
		assert.NotNil(t, conn)
	}
}

func TestVirtualListenerAcceptAfterClose(t *testing.T) {
	vl := newVirtualListener(testAddr{})
	vl.Close()

	_, err := vl.Accept()
	assert.Equal(t, ErrListenerClosed, err)

	_, err = vl.Accept()
	assert.Equal(t, ErrListenerClosed, err)
}
