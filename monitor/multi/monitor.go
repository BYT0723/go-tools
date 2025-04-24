package multi

import (
	"context"

	"github.com/BYT0723/go-tools/monitor"
	"github.com/BYT0723/go-tools/monitor/component"
)

type (
	// MultiMonitor is a struct that combines multiple Monitor instances.
	// It manages multiple monitors and their alerts, allowing concurrent monitoring.
	MultiMonitor struct {
		*component.MonitorComponent                   // Embedded MonitorComponent to handle alert notifications and monitor lifecycle.
		monitors                    []monitor.Monitor // A slice of monitors to be managed.
	}
)

// NewMultiMonitor creates and returns a new MultiMonitor instance with the provided monitors.
// monitors: A list of Monitor instances that will be managed and monitored concurrently.
func NewMultiMonitor(monitors ...monitor.Monitor) *MultiMonitor {
	return &MultiMonitor{
		MonitorComponent: component.NewMonitorComponent(), // Initialize MonitorComponent with default values.
		monitors:         monitors,                        // Set the provided monitors.
	}
}

// Start begins the monitoring of multiple monitors concurrently. Each monitor is started in its own goroutine,
// and alerts from all monitors are forwarded to the MultiMonitor's alert channel.
func (m *MultiMonitor) Start(ctx context.Context) {
	m.SetContext(ctx) // Set the context for the MultiMonitor.

	// Start each monitor concurrently.
	for _, item := range m.monitors {
		go func(monitorInstance monitor.Monitor) {
			monitorInstance.Start(m.Context()) // Start the individual monitor.
			// Forward alerts from the individual monitor to the MultiMonitor's alert channel.
			for ar := range monitorInstance.Subscribe() {
				m.Notify(ar)
			}
		}(item)
	}

	// Start a goroutine to listen for the context cancellation and stop the MultiMonitor when done.
	go func() {
		<-m.Context().Done()
		m.Stop(ctx) // Stop all monitors and close resources when the context is cancelled.
	}()
}
