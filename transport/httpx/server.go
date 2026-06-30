package httpx

import (
	"context"
	"net"
	"net/http"

	"github.com/BYT0723/go-tools/transport/connmux"
)

type httpService struct {
	srv *http.Server
	l   net.Listener
}

func WrapHTTPServer(srv *http.Server) connmux.ListenedService {
	return &httpService{srv: srv}
}

func (hs *httpService) Name() string                      { return "http" }
func (hs *httpService) Match() connmux.Matcher             { return connmux.MatchHTTP1 }
func (hs *httpService) SetListener(l net.Listener)         { hs.l = l }
func (hs *httpService) Init(_ context.Context) error       { return nil }
func (hs *httpService) Run(_ context.Context) error {
	if err := hs.srv.Serve(hs.l); err != http.ErrServerClosed {
		return err
	}
	return nil
}
func (hs *httpService) Destroy(_ context.Context) error { return hs.srv.Close() }
