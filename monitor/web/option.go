package web

import (
	"time"
)

func WithTimeout(timeout time.Duration) Option {
	return func(p *Monitor) {
		p.Timeout = timeout
	}
}
