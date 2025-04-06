package component

import (
	"context"
	"time"

	"github.com/BYT0723/go-tools/channelx"
	"github.com/BYT0723/go-tools/monitor"
)

type MonitorComponent struct {
	ctx context.Context
	cf  context.CancelFunc

	cycle       time.Duration
	cycleTicker *time.Ticker

	Timeout time.Duration

	ch chan *monitor.Alert
}

func NewMonitorComponent() *MonitorComponent {
	return &MonitorComponent{
		cycle:   time.Minute,
		Timeout: time.Minute,
		ch:      make(chan *monitor.Alert, 1024),
	}
}

func (m *MonitorComponent) SetCycle(cycle time.Duration) {
	if m.cycle == cycle || cycle == 0 {
		return
	}
	m.cycle = cycle
	if m.Timeout > m.cycle {
		m.Timeout = m.cycle
	}
	if m.cycleTicker != nil {
		m.cycleTicker.Reset(cycle)
	}
}

func (m *MonitorComponent) SetTimeout(timeout time.Duration) {
	if m.Timeout == timeout || timeout == 0 {
		return
	}
	m.Timeout = timeout
}

func (m *MonitorComponent) Notify(as ...*monitor.Alert) {
	for _, a := range as {
		channelx.ChannelInWithTimeout(m.ch, a, time.Second)
	}
}

func (m *MonitorComponent) Subscribe() <-chan *monitor.Alert {
	return m.ch
}

func (m *MonitorComponent) Stop(ctx context.Context) {
	close(m.ch)
	if m.cf != nil {
		m.cf()
	}
}

func (m *MonitorComponent) SetContext(ctx context.Context) {
	m.ctx, m.cf = context.WithCancel(ctx)
}

func (m *MonitorComponent) Context() context.Context {
	if m.ctx == nil {
		m.ctx, m.cf = context.WithCancel(context.Background())
	}
	return m.ctx
}

func (m *MonitorComponent) Ticker() *time.Ticker {
	if m.cycleTicker == nil {
		m.cycleTicker = time.NewTicker(m.cycle)
	}
	return m.cycleTicker
}
