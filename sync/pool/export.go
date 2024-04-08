package pool

import (
	"context"
	"sync"
)

type Pool struct{}

type Group struct {
	wg   sync.WaitGroup
	jobs []*Job
}

func (g *Group) AddJob(j *Job) {
	g.jobs = append(g.jobs, j)
}

func (g *Group) Run(pctx context.Context) {
	for _, job := range g.jobs {
		g.wg.Add(1)
		go func(pctx context.Context, wg *sync.WaitGroup, job *Job) {
			defer wg.Done()
		out:
			for {
				select {
				case <-pctx.Done():
					break out
				default:
					job.cond.L.Lock()
					for job.cond.status == PAUSE {
						job.cond.Wait()
					}
					if job.cond.status == EXIT {
						job.cond.L.Unlock()
						break out
					}
					job.cond.L.Unlock()
					job.run(pctx)
				}
			}
		}(pctx, &g.wg, job)
	}
}

func (g *Group) Wait() {
	g.wg.Wait()
}
