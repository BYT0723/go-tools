package srvx

import (
	"context"
	"errors"
	"testing"

	"github.com/BYT0723/go-tools/logx"
	"github.com/BYT0723/go-tools/logx/logcore"
)

type (
	HTTPService struct{}
	SSHService  struct{}
	FRPService  struct{}
)

func (s *HTTPService) Name() string { return "http" }
func (s *SSHService) Name() string  { return "ssh" }
func (s *FRPService) Name() string  { return "frp" }

func (s *HTTPService) Start(ctx context.Context) error {
	return errors.New("未实现的服务")
}

func (s *SSHService) Start(ctx context.Context) error {
	return nil
}

func (s *FRPService) Start(ctx context.Context) error {
	return nil
}

func (s *HTTPService) ErrHandle(ctx context.Context, err error) {}
func (s *SSHService) ErrHandle(ctx context.Context, err error)  {}
func (s *FRPService) ErrHandle(ctx context.Context, err error)  {}

func TestServices_Run(t *testing.T) {
	if err := logx.Init(logx.WithConf(&logcore.LoggerConf{Console: true})); err != nil {
		t.Fatal(err)
	}
	t.Run("test", func(t *testing.T) {
		var srv Services
		srv.Log = logx.Default()

		srv.Register(&HTTPService{})
		srv.Register(&SSHService{})
		srv.Register(&FRPService{})

		srv.Run(context.Background())
	})
}
