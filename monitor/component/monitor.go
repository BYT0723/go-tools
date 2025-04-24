// Package component provides a monitor component that handles periodic monitoring
// with alert notifications. It allows setting the cycle time, timeout, and subscribing to alerts.
package component

import (
	"context"
	"time"

	"github.com/BYT0723/go-tools/channelx"
	"github.com/BYT0723/go-tools/monitor"
)

type (
	// MonitorComponent is a struct that encapsulates monitoring functionality.
	// It manages the monitoring cycle, timeout, and alert notifications.
	MonitorComponent struct {
		ctx         context.Context     // The context for cancellation
		cf          context.CancelFunc  // The cancel function for the context
		cycle       time.Duration       // The monitoring cycle duration
		cycleTicker *time.Ticker        // The ticker for the monitoring cycle
		Timeout     time.Duration       // The timeout duration for monitoring
		ch          chan *monitor.Alert // Channel for alerts
	}
)

// NewMonitorComponent creates and returns a new MonitorComponent instance with default settings.
func NewMonitorComponent() *MonitorComponent {
	return &MonitorComponent{
		cycle:   time.Minute,                     // Default cycle is 1 minute
		Timeout: time.Minute,                     // Default timeout is 1 minute
		ch:      make(chan *monitor.Alert, 1024), // Create a channel for alerts with a buffer size of 1024
	}
}

// SetCycle sets the monitoring cycle duration. If the cycle is 0 or unchanged, it will not update.
// If the timeout is greater than the new cycle, the timeout will be adjusted to match the new cycle.
func (m *MonitorComponent) SetCycle(cycle time.Duration) {
	if m.cycle == cycle || cycle == 0 {
		return
	}
	m.cycle = cycle
	if m.Timeout > m.cycle {
		m.Timeout = m.cycle
	}
	if m.cycleTicker != nil {
		m.cycleTicker.Reset(cycle) // Reset the ticker with the new cycle duration
	}
}

// SetTimeout sets the timeout duration for monitoring. If the timeout is 0 or unchanged, it will not update.
func (m *MonitorComponent) SetTimeout(timeout time.Duration) {
	if m.Timeout == timeout || timeout == 0 {
		return
	}
	m.Timeout = timeout
}

// Notify sends one or more alerts to the alert channel. It uses a timeout when sending the alert.
func (m *MonitorComponent) Notify(as ...*monitor.Alert) {
	for _, a := range as {
		// Send the alert to the channel with a timeout of 1 second
		_ = channelx.InTimeout(m.ch, a, time.Second)
	}
}

// Subscribe returns a read-only channel to subscribe to alerts.
func (m *MonitorComponent) Subscribe() <-chan *monitor.Alert {
	return m.ch
}

// Stop stops the monitoring and closes the alert channel. It also cancels the context.
func (m *MonitorComponent) Stop(ctx context.Context) {
	close(m.ch) // Close the alert channel
	if m.cf != nil {
		m.cf() // Cancel the context if it exists
	}
}

// SetContext sets the context for the MonitorComponent. It creates a new context with cancel function.
func (m *MonitorComponent) SetContext(ctx context.Context) {
	m.ctx, m.cf = context.WithCancel(ctx) // Create a cancelable context
}

// Context returns the current context for the MonitorComponent. If the context is nil, it initializes it.
func (m *MonitorComponent) Context() context.Context {
	if m.ctx == nil {
		m.ctx, m.cf = context.WithCancel(context.Background()) // Initialize the context if it is nil
	}
	return m.ctx
}

// Ticker returns the current ticker for the monitoring cycle. If the ticker is nil, it initializes it.
func (m *MonitorComponent) Ticker() *time.Ticker {
	if m.cycleTicker == nil {
		m.cycleTicker = time.NewTicker(m.cycle) // Create a new ticker if not already initialized
	}
	return m.cycleTicker
}
