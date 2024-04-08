package pool

import (
	"context"
)

type (
	RunFunc func(context.Context)
)

type Job struct {
	ctx  context.Context
	cond *Cond
	run  RunFunc
}

func NewJob(run RunFunc) *Job {
	return &Job{
		run: run,
	}
}
