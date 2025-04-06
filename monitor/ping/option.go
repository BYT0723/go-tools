package ping

import "time"

func WithCount(count int) Option {
	return func(p *Monitor) {
		p.count = count
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(p *Monitor) {
		p.Timeout = timeout
	}
}

func WithICMP() Option {
	return func(p *Monitor) {
		p.icmp = true
	}
}
