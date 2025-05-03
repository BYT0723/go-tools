package prometheus

import "time"

type Option func(*Monitor)

func WithCycle(cycle time.Duration) Option {
	return func(m *Monitor) {
		m.SetCycle(cycle)
	}
}

func WithTimeout(t time.Duration) Option {
	return func(m *Monitor) {
		m.Timeout = t
	}
}
