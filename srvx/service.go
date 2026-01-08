package srvx

import (
	"context"
)

type Service interface {
	Run(ctx context.Context) error
	ErrHandle(ctx context.Context, err error)
}

type Services []Service

func (s Services) Run(ctx context.Context) {
	for _, s := range s {
		if err := s.Run(ctx); err != nil {
			s.ErrHandle(ctx, err)
		}
	}
}
