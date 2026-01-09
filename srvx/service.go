package srvx

import (
	"context"
	"sync"

	ctxx "github.com/BYT0723/go-tools/contextx"
	"github.com/BYT0723/go-tools/logx"
)

type Service interface {
	Name() string
	// 非阻塞式运行
	// 可以从ctx中获取WaitGroup去启动goroutine防止阻塞
	Start(ctx context.Context) error
	ErrHandle(ctx context.Context, err error)
}

type Services struct {
	wg       sync.WaitGroup
	services []Service
	Log      logx.Logger
}

func (ss *Services) Register(s Service) {
	ss.services = append(ss.services, s)
}

// Run start all service and wait for all service exit
func (ss *Services) Run(ctx context.Context) {
	ctx = ctxx.WithWaitGroup(ctx, &ss.wg)
	for _, s := range ss.services {
		if err := s.Start(ctx); err != nil {
			s.ErrHandle(ctx, err)
		}
		if ss.Log != nil {
			ss.Log.Infof("service %s run success", s.Name())
		}
	}
}
