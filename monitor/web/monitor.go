package web

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/BYT0723/go-tools/monitor"
	"github.com/BYT0723/go-tools/monitor/component"
	"github.com/BYT0723/go-tools/transport/httpx"
)

var _ monitor.Monitor = (*Monitor)(nil)

type (
	// Monitor is a struct that monitors a web service by sending HTTP requests and evaluating the responses.
	// It embeds MonitorComponent and AlertComponent to handle monitoring and alerting.
	Monitor struct {
		*component.MonitorComponent              // Embedding MonitorComponent to handle lifecycle and alert notification.
		component.AlertComponent[httpx.Response] // Embedding AlertComponent to handle alert rules based on HTTP responses.

		client  *httpx.Client // The HTTP client used for sending requests.
		method  string        // The HTTP method (GET, POST, etc.).
		addr    string        // The target address of the web service.
		header  http.Header   // Custom headers for the HTTP request.
		payload any           // The request payload (body), if applicable.
	}

	// Option is a function type for configuring the Monitor.
	Option func(*Monitor)
)

var (
	// Default HTTP client configuration with TLS settings and redirect handling.
	defaultClient = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Disable TLS verification for testing.
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Prevent more than 10 redirects.
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
	// Default monitoring cycle (interval between HTTP requests).
	defautlCycle = time.Minute
)

// NewMonitor creates a new Monitor instance that monitors a web service via HTTP requests.
// method: The HTTP method (e.g., GET, POST).
// addr: The target URL of the web service.
// opts: Options for configuring the monitor (e.g., headers, payload).
func NewMonitor(method string, addr string, opts ...Option) *Monitor {
	p := &Monitor{
		MonitorComponent: component.NewMonitorComponent(),                       // Initialize MonitorComponent with default values.
		client:           httpx.NewClient(httpx.WithHttpClient(&defaultClient)), // Initialize the HTTP client with the default configuration.
		addr:             addr,                                                  // Set the target address for the web service.
	}
	p.SetCycle(defautlCycle) // Set the default monitoring cycle.
	p.Timeout = defautlCycle // Set the default timeout.

	// Apply any provided options to configure the monitor.
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Start begins the web monitoring process, sending HTTP requests at the specified cycle interval.
// It listens for alerts and evaluates the responses based on the configured alert rules.
func (m *Monitor) Start(ctx context.Context) {
	m.SetContext(ctx) // Set the context for the monitor.

	// Start the monitoring process in a separate goroutine.
	go func() {
		t := m.Ticker() // Get the ticker for the monitoring cycle.
		defer t.Stop()

		// Perform the first HTTP request and evaluate the response.
		s, err := m.do()
		if err != nil {
			m.Notify(monitor.InternalAlert(err)) // Notify of an internal error if the request fails.
		} else {
			m.Notify(m.Evaluate(s)...)
		}

		// Continuously monitor the web service by sending HTTP requests at the specified cycle interval.
		for {
			select {
			case <-m.Context().Done():
				m.Stop(ctx) // Stop monitoring if the context is canceled.
				return
			case <-t.C:
				// Perform subsequent HTTP requests at regular intervals.
				s, err := m.do()
				if err != nil {
					m.Notify(monitor.InternalAlert(err)) // Notify error if the request fails.
					continue
				}

				m.Notify(m.Evaluate(s)...) // Notify the alerts based on the HTTP response.
			}
		}
	}()
}

// do performs the actual HTTP request to the target address and returns the response or an error.
// It uses a context with timeout for each request.
func (m *Monitor) do() (*httpx.Response, error) {
	ctx, cf := context.WithTimeout(context.Background(), m.Timeout) // Create a context with timeout.
	defer cf()                                                      // Ensure the cancel function is called.

	// Send the HTTP request using the configured client, method, address, headers, and payload.
	return m.client.Do(
		ctx,
		m.method,
		m.addr,
		httpx.WithHeader(m.header),
		httpx.WithPayload(m.payload),
	)
}
