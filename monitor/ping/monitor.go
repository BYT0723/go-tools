package ping

import (
	"time"

	"github.com/BYT0723/go-tools/monitor"
	"github.com/BYT0723/go-tools/monitor/component"
	probing "github.com/prometheus-community/pro-bing"
	"golang.org/x/net/context"
)

var _ monitor.Monitor = (*Monitor)(nil)

type (
	Monitor struct {
		component.MonitorComponent
		component.AlertComponent[probing.Statistics]

		addr  string
		count int
		icmp  bool
	}
	Option func(p *Monitor)
)

var (
	cycle    = time.Minute
	count    = 10
	minCount = 3
)

func NewMonitor(addr string, opts ...Option) *Monitor {
	p := &Monitor{
		MonitorComponent: *component.NewMonitorComponent(),
		addr:             addr,
		count:            count,
	}
	p.SetCycle(cycle)
	for _, opt := range opts {
		opt(p)
	}
	p.count = max(minCount, p.count)
	// p.count * p.interval (interval default 1s)
	p.Timeout = max(p.Timeout, time.Duration(p.count)*time.Second)

	return p
}

func (m *Monitor) Start(ctx context.Context) {
	m.SetContext(ctx)

	go func() {
		t := m.MonitorComponent.Ticker()
		defer t.Stop()

		s, err := m.do()
		if err != nil {
			m.Notify(monitor.InternalAlert(err))
		} else {
			m.Notify(m.Evaluate(s)...)
		}

		for {
			select {
			case <-m.Context().Done():
				m.MonitorComponent.Stop(ctx)
				return
			case <-t.C:
				s, err := m.do()
				if err != nil {
					m.Notify(monitor.InternalAlert(err))
					continue
				}

				m.Notify(m.Evaluate(s)...)
			}
		}
	}()
}

func (m *Monitor) do() (stats *probing.Statistics, err error) {
	pinger, err := probing.NewPinger(m.addr)
	if err != nil {
		return
	}
	pinger.Count = m.count
	pinger.Timeout = m.Timeout

	if m.icmp {
		pinger.SetPrivileged(true)
	}

	if err = pinger.Run(); err != nil {
		return
	}
	stats = pinger.Statistics()
	return
}
