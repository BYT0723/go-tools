package connmux

import (
	"fmt"
	"net"
	"sync"
)

var ErrListenerClosed = fmt.Errorf("virtual listener closed")

type VirtualListener struct {
	ch     chan net.Conn
	closed chan struct{}
	addr   net.Addr
	mu     sync.Mutex
}

const virtualListenerBuffer = 128

func newVirtualListener(addr net.Addr) *VirtualListener {
	return &VirtualListener{
		ch:     make(chan net.Conn, virtualListenerBuffer),
		closed: make(chan struct{}),
		addr:   addr,
	}
}

func (vl *VirtualListener) Accept() (net.Conn, error) {
	select {
	case conn := <-vl.ch:
		return conn, nil
	case <-vl.closed:
		return nil, ErrListenerClosed
	}
}

func (vl *VirtualListener) Close() error {
	vl.mu.Lock()
	defer vl.mu.Unlock()
	select {
	case <-vl.closed:
		return nil
	default:
		close(vl.closed)
	}
	return nil
}

func (vl *VirtualListener) Addr() net.Addr {
	return vl.addr
}

func (vl *VirtualListener) push(conn net.Conn) {
	select {
	case vl.ch <- conn:
	case <-vl.closed:
		conn.Close()
	default:
		conn.Close()
	}
}
