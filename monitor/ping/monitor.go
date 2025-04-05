package ping

import (
	"sync"
	"time"

	"github.com/BYT0723/go-tools/monitor"
	probing "github.com/prometheus-community/pro-bing"
	"golang.org/x/net/context"
)

var _ monitor.Monitor = (*Monitor)(nil)

type (
	Monitor struct {
		addr            string
		cycle           time.Duration
		count           int
		alert           bool
		icmp            bool
		timeout         time.Duration // 单次ping操作的超时时间
		alertRulesMutex sync.Mutex
		alertRules      []monitor.AlertRule[probing.Statistics]
		ch              chan *monitor.Alert
		ctx             context.Context
		cf              context.CancelFunc
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
		addr:  addr,
		cycle: cycle,
		count: count,
		ch:    make(chan *monitor.Alert, 1024),
	}
	for _, opt := range opts {
		opt(p)
	}
	p.count = max(minCount, p.count)
	// p.count * p.interval (interval default 1s)
	p.timeout = max(p.timeout, time.Duration(p.count)*time.Second)
	p.cycle = max(p.timeout, p.cycle)

	return p
}

func (m *Monitor) Start(ctx context.Context) {
	m.ctx, m.cf = context.WithCancel(ctx)

	go func() {
		ticker := time.NewTicker(m.cycle)
		defer ticker.Stop()

		s, err := m.do()
		if err != nil {
			m.ch <- monitor.InternalAlert(err)
		} else {
			m.alertRulesMutex.Lock()
			for _, ar := range m.alertRules {
				if a, b := ar(s); b {
					m.ch <- a
				}
			}
			m.alertRulesMutex.Unlock()
		}

		for {
			select {
			case <-m.ctx.Done():
				m.cf()
				close(m.ch)
				return
			case <-ticker.C:
				s, err := m.do()
				if err != nil {
					m.ch <- monitor.InternalAlert(err)
					continue
				}

				for _, ar := range m.alertRules {
					if a, b := ar(s); b {
						m.ch <- a
					}
				}
			}
		}
	}()
}

func (m *Monitor) Stop(ctx context.Context) {
	close(m.ch)
	if m.cf != nil {
		m.cf()
	}
}

func (m *Monitor) Subscribe() <-chan *monitor.Alert {
	return m.ch
}

func (m *Monitor) AddAlertRule(ars ...monitor.AlertRule[probing.Statistics]) {
	m.alertRulesMutex.Lock()
	m.alertRules = append(m.alertRules, ars...)
	m.alertRulesMutex.Unlock()
}

func (m *Monitor) do() (stats *probing.Statistics, err error) {
	pinger, err := probing.NewPinger(m.addr)
	if err != nil {
		return
	}
	pinger.Count = m.count
	pinger.Timeout = m.timeout

	if m.icmp {
		pinger.SetPrivileged(true)
	}

	if err = pinger.Run(); err != nil {
		return
	}
	stats = pinger.Statistics()
	return
}
