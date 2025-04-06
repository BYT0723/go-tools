package multi

import (
	"context"

	"github.com/BYT0723/go-tools/monitor"
	"github.com/BYT0723/go-tools/monitor/component"
)

type MultiMonitor struct {
	*component.MonitorComponent
	monitors []monitor.Monitor
}

func NewMultiMonitor(monitors ...monitor.Monitor) *MultiMonitor {
	return &MultiMonitor{
		MonitorComponent: component.NewMonitorComponent(),
		monitors:         monitors,
	}
}

func (m *MultiMonitor) Start(ctx context.Context) {
	m.SetContext(ctx)

	for _, item := range m.monitors {
		go func() {
			item.Start(m.Context())
			for ar := range item.Subscribe() {
				m.Notify(ar)
			}
		}()
	}

	go func() {
		<-m.Context().Done()
		m.Stop(ctx)
	}()
}
