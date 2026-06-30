package connmux

import "time"

type Option func(*Mux)

func WithAddr(addr string) Option {
	return func(m *Mux) { m.addr = addr }
}

func WithSniffSize(n int) Option {
	return func(m *Mux) { m.sniffSize = n }
}

func WithSniffDeadline(d time.Duration) Option {
	return func(m *Mux) { m.sniffDeadline = d }
}
