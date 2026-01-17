package srvx

import (
	"context"
	"sync"

	"github.com/BYT0723/go-tools/logx"
	"github.com/BYT0723/go-tools/logx/noplogger"
)

type Service interface {
	Name() string

	// Init 初始化资源
	// 返回 error 表示服务无法启动
	Init(ctx context.Context) error

	// Run 启动后台任务
	// Run 只在 Init 成功后调用
	// 返回 error 表示服务提前异常退出
	Run(ctx context.Context) error

	// Destroy 释放资源
	// 只会在 Init 成功后被调用
	// 必须幂等
	Destroy(ctx context.Context) error
}

type Services struct {
	Log logx.Logger

	wg       sync.WaitGroup
	services []Service
}

func (ss *Services) Register(s Service) {
	ss.services = append(ss.services, s)
}

// Run start all service and wait for all service exit
func (ss *Services) Run(ctx context.Context) {
	if ss.Log == nil {
		ss.Log = noplogger.NopLogger{}
	}
	for _, s := range ss.services {
		ss.wg.Go(func() {
			name := s.Name()

			ss.Log.Info("service init", logx.String("name", name))
			if err := s.Init(ctx); err != nil {
				ss.Log.Error("service init error", logx.String("name", name), logx.Err(err))
				return
			}

			defer func() {
				ss.Log.Error("service exit", logx.String("name", name))
				if err := s.Destroy(ctx); err != nil {
					ss.Log.Error("service destroy error", logx.String("name", name), logx.Err(err))
				}
			}()

			ss.Log.Info("service run", logx.String("name", name))
			if err := s.Run(ctx); err != nil {
				ss.Log.Error("service run error", logx.String("name", name), logx.Err(err))
				return
			}
		})
	}
	ss.wg.Wait()
}
