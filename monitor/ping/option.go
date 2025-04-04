package ping

import (
	"time"

	"github.com/BYT0723/go-tools/monitor"
	probing "github.com/prometheus-community/pro-bing"
)

func WithCycle(cycle time.Duration) Option {
	return func(p *Monitor) {
		p.cycle = cycle
	}
}

func WithCount(count int) Option {
	return func(p *Monitor) {
		p.count = count
	}
}

func WithAlertRules(rules ...monitor.AlertRule[probing.Statistics]) Option {
	return func(p *Monitor) {
		p.alertRules = rules
	}
}

func WithICMP() Option {
	return func(p *Monitor) {
		p.icmp = true
	}
}
