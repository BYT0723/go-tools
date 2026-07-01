package tlsx

import "net"

type SingleConnListener struct {
	Conn net.Conn
	done chan struct{}
}

func (l *SingleConnListener) Accept() (net.Conn, error) {
	if l.done == nil {
		l.done = make(chan struct{})
		return l.Conn, nil
	}
	<-l.done
	return nil, net.ErrClosed
}

func (l *SingleConnListener) Close() error {
	if l.done != nil {
		close(l.done)
	}
	return l.Conn.Close()
}

func (l *SingleConnListener) Addr() net.Addr { return l.Conn.LocalAddr() }
