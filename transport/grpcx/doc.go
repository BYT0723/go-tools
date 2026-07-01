// Package grpcx provides gRPC utilities for Go services.
//
// # Server
//
// WrapGRPCServer wraps a *grpc.Server into connmux.ListenedService,
// allowing it to be multiplexed with other protocols on a single port:
//
//	grpcSrv := grpc.NewServer()
//	pb.RegisterEchoServer(grpcSrv, &echoImpl{})
//
//	mux := connmux.NewMux(connmux.WithAddr(":443"))
//	mux.Route("grpc", grpcx.WrapGRPCServer(grpcSrv))
package grpcx
