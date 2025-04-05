package multi

import (
	"context"

	"github.com/BYT0723/go-tools/monitor"
)

type MultiMonitor struct {
	ch       chan *monitor.Alert
	monitors []monitor.Monitor
	ctx      context.Context
	cf       context.CancelFunc
}

func NewMultiMonitor(monitors ...monitor.Monitor) *MultiMonitor {
	return &MultiMonitor{
		monitors: monitors,
		ch:       make(chan *monitor.Alert, 1024*len(monitors)),
	}
}

func (m *MultiMonitor) Start(ctx context.Context) {
	m.ctx, m.cf = context.WithCancel(ctx)
	for _, item := range m.monitors {
		go func() {
			item.Start(m.ctx)
			for ar := range item.Subscribe() {
				m.ch <- ar
			}
		}()
	}
}

func (m *MultiMonitor) Stop(ctx context.Context) {
	m.cf()
	close(m.ch)
}

func (m *MultiMonitor) Subscribe() <-chan *monitor.Alert {
	return m.ch
}
