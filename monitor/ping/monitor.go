package ping

import (
	"time"

	"github.com/BYT0723/go-tools/monitor"
	"github.com/BYT0723/go-tools/monitor/component"
	probing "github.com/prometheus-community/pro-bing"
	"golang.org/x/net/context"
)

var _ monitor.Monitor = (*Monitor)(nil)

type (
	// Monitor is a struct that performs ping monitoring on a given network address.
	// It can monitor the availability of a host using ICMP (Ping) or other methods.
	Monitor struct {
		*component.MonitorComponent                  // Embedding MonitorComponent to handle lifecycle and alert notification.
		component.AlertComponent[probing.Statistics] // Embedding AlertComponent to handle alert rules for ping statistics.

		addr  string // The target address to ping.
		count int    // The number of ping attempts to make.
		icmp  bool   // Whether to use ICMP (Ping) or other methods.
	}

	// Option is a function type for configuring the Monitor.
	Option func(p *Monitor)
)

var (
	// Default cycle for monitoring.
	cycle = time.Minute
	// Default count of ping attempts per monitoring cycle.
	count = 10
	// Minimum count of packets required for the monitoring cycle.
	minCount = 3
)

// NewMonitor creates a new Monitor instance for pinging a specific address with the provided options.
// addr: The target address to ping.
// opts: Options for configuring the monitor such as count, timeout, etc.
func NewMonitor(addr string, opts ...Option) *Monitor {
	p := &Monitor{
		MonitorComponent: component.NewMonitorComponent(), // Initialize MonitorComponent with default values.
		addr:             addr,                            // Set the target address.
		count:            count,                           // Set the default ping attempt count.
	}
	p.SetCycle(cycle) // Set the monitoring cycle duration.

	// Apply any provided options to configure the monitor.
	for _, opt := range opts {
		opt(p)
	}
	// Ensure the ping attempt count is above the minimum threshold.
	p.count = max(minCount, p.count)
	// Adjust timeout based on the count and default interval.
	p.Timeout = max(p.Timeout, time.Duration(p.count)*time.Second)

	return p
}

// Start begins the ping monitoring process, where it pings the target address at the specified cycle.
// It listens for alerts and pings the address regularly based on the defined cycle.
func (m *Monitor) Start(ctx context.Context) {
	m.SetContext(ctx) // Set the context for the monitor.

	// Start the monitoring process in a separate goroutine.
	go func() {
		t := m.Ticker() // Get the ticker for the monitoring cycle.
		defer t.Stop()

		// Perform the first ping operation.
		s, err := m.do()
		if err != nil {
			m.Notify(monitor.InternalAlert(err)) // Notify of an internal error if the ping fails.
		} else {
			m.Notify(m.Evaluate(s)...) // Notify the alerts based on the ping statistics.
		}

		// Continuously ping the address at the specified cycle interval.
		for {
			select {
			case <-m.Context().Done():
				m.Stop(ctx) // Stop monitoring if the context is canceled.
				return
			case <-t.C:
				// Perform subsequent pings at regular intervals.
				s, err := m.do()
				if err != nil {
					m.Notify(monitor.InternalAlert(err)) // Notify error if ping fails.
					continue
				}

				m.Notify(m.Evaluate(s)...) // Notify the alerts based on the ping statistics.
			}
		}
	}()
}

// do performs the actual ping operation, using the specified address and configuration.
// It returns the ping statistics or an error if the operation fails.
func (m *Monitor) do() (stats *probing.Statistics, err error) {
	pinger, err := probing.NewPinger(m.addr) // Create a new Pinger instance with the target address.
	if err != nil {
		return
	}
	pinger.Count = m.count     // Set the number of ping attempts.
	pinger.Timeout = m.Timeout // Set the timeout for the ping.

	// If ICMP is enabled, set the pinger to use privileged mode.
	if m.icmp {
		pinger.SetPrivileged(true)
	}

	// Run the ping operation.
	if err = pinger.Run(); err != nil {
		return
	}
	stats = pinger.Statistics() // Retrieve the statistics from the ping operation.
	return
}
