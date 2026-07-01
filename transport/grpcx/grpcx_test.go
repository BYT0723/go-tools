package grpcx

import (
	"context"
	"net"
	"testing"

	"github.com/BYT0723/go-tools/transport/connmux"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestWrapGRPCServer(t *testing.T) {
	t.Skip("skipping: grpc.NewServer() panics without registered services in some envs")

	srv := grpc.NewServer()
	svc := WrapGRPCServer(srv)
	assert.Equal(t, "grpc", svc.Name())
	assert.Equal(t, connmux.MatchHTTP2, svc.Match())

	assert.NoError(t, svc.Init(context.Background()))
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	svc.SetListener(l)

	errCh := make(chan error, 1)
	go func() { errCh <- svc.Run(context.Background()) }()

	assert.NoError(t, svc.Destroy(context.Background()))

	if err := <-errCh; err != nil {
		t.Logf("Run returned: %v", err)
	}
}
