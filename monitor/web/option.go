package web

import (
	"time"

	"github.com/BYT0723/go-tools/monitor"
)

func WithCycle(cycle time.Duration) Option {
	return func(p *Monitor) {
		p.cycle = cycle
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(p *Monitor) {
		p.timeout = timeout
	}
}

func WithAlertRules(rules ...monitor.AlertRule[Statistics]) Option {
	return func(p *Monitor) {
		p.alertRules = rules
	}
}
