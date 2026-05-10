package srvx

import (
	"context"
	"testing"

	"github.com/BYT0723/go-tools/logx"
	"github.com/BYT0723/go-tools/logx/logcore"
	"github.com/BYT0723/go-tools/logx/noplogger"
)

type (
	httpService struct{}
	sshService  struct{}
	frpService  struct{}
)

func (s *httpService) Name() string { return "http" }
func (s *sshService) Name() string  { return "ssh" }
func (s *frpService) Name() string  { return "frp" }

func (s *httpService) Init(ctx context.Context) error  { return nil }
func (s *sshService) Init(ctx context.Context) error   { return nil }
func (s *frpService) Init(ctx context.Context) error   { return nil }

func (s *httpService) Run(ctx context.Context) error   { return nil }
func (s *sshService) Run(ctx context.Context) error    { return nil }
func (s *frpService) Run(ctx context.Context) error    { return nil }

func (s *httpService) Destroy(ctx context.Context) error { return nil }
func (s *sshService) Destroy(ctx context.Context) error  { return nil }
func (s *frpService) Destroy(ctx context.Context) error  { return nil }

func TestServices_Run(t *testing.T) {
	if err := logx.Init(logx.WithConf(&logcore.LoggerConf{Console: true})); err != nil {
		t.Fatal(err)
	}

	var srv Services
	srv.Log = logx.Default()

	srv.Register(&httpService{})
	srv.Register(&sshService{})
	srv.Register(&frpService{})

	srv.Run(context.Background())
}

func TestServices_NopLogger(t *testing.T) {
	var srv Services

	srv.Register(&httpService{})
	srv.Register(&sshService{})

	srv.Run(context.Background())
}

func TestServices_InitError(t *testing.T) {
	type initErrService struct {
		httpService
	}
	// 暂时不测试 init error 的完整路径，因为 sync.WaitGroup 不支持 goroutine 管理
}

func TestServiceInterface(t *testing.T) {
	var _ Service = (*httpService)(nil)
	var _ Service = (*sshService)(nil)
	var _ Service = (*frpService)(nil)

	// 编译期验证 NopLogger 可作为 Logger
	var _ logx.Logger = noplogger.NopLogger{}
}
