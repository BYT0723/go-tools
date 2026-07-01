package grpcx

import (
	"context"
	"net"

	"github.com/BYT0723/go-tools/transport/connmux"
	"google.golang.org/grpc"
)

type grpcService struct {
	srv *grpc.Server
	l   net.Listener
}

func WrapGRPCServer(srv *grpc.Server) connmux.ListenedService {
	return &grpcService{srv: srv}
}

func (gs *grpcService) Name() string                   { return "grpc" }
func (gs *grpcService) Match() connmux.Matcher          { return connmux.MatchHTTP2 }
func (gs *grpcService) SetListener(l net.Listener)      { gs.l = l }
func (gs *grpcService) Init(_ context.Context) error    { return nil }
func (gs *grpcService) Run(_ context.Context) error     { return gs.srv.Serve(gs.l) }
func (gs *grpcService) Destroy(_ context.Context) error { gs.srv.Stop(); return nil }
