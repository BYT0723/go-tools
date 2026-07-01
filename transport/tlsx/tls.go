package tlsx

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/BYT0723/go-tools/transport/connmux"
)

type tlsService struct {
	tlsCfg  *tls.Config
	handler func(net.Conn)
	l       net.Listener
}

func WrapTLS(tlsCfg *tls.Config, handler func(net.Conn)) connmux.ListenedService {
	return &tlsService{tlsCfg: tlsCfg, handler: handler}
}

func (ts *tlsService) Name() string                   { return "tls" }
func (ts *tlsService) Match() connmux.Matcher          { return connmux.MatchTLS }
func (ts *tlsService) SetListener(l net.Listener)       { ts.l = l }
func (ts *tlsService) Init(_ context.Context) error     { return nil }
func (ts *tlsService) Destroy(_ context.Context) error { return ts.l.Close() }
func (ts *tlsService) Run(_ context.Context) error {
	for {
		conn, err := ts.l.Accept()
		if err != nil {
			return err
		}
		go func() {
			tlsConn := tls.Server(conn, ts.tlsCfg)
			defer tlsConn.Close()
			if err := tlsConn.Handshake(); err != nil {
				return
			}
			ts.handler(tlsConn)
		}()
	}
}
