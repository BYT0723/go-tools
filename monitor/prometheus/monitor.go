package prometheus

import (
	"bytes"
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/BYT0723/go-tools/monitor"
	"github.com/BYT0723/go-tools/monitor/component"
	"github.com/BYT0723/go-tools/transport/httpx"
	promModel "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

var _ monitor.Monitor = (*Monitor)(nil)

type (
	MetricFamily = promModel.MetricFamily
	Monitor      struct {
		*component.MonitorComponent
		component.AlertComponent[Data]

		client *httpx.Client
		addr   string
	}
	Data struct {
		Metric map[string]*MetricFamily
	}
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

func NewMonitor(addr string, opts ...Option) *Monitor {
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

// Start implements monitor.Monitor.
func (m *Monitor) Start(ctx context.Context) {
	m.SetContext(ctx)

	go func() {
		t := m.Ticker()
		defer t.Stop()

		for {
			select {
			case <-m.Context().Done():
				m.Stop(ctx)
				return
			case <-t.C:
				d, err := m.do()
				if err != nil {
					m.Notify(monitor.InternalAlert(err))
					continue
				}
				m.Notify(m.Evaluate(d)...)
			}
		}
	}()
}

func (m *Monitor) do() (*Data, error) {
	ctx, cf := context.WithTimeout(m.Context(), m.Timeout)
	defer cf()

	resp, err := m.client.Get(ctx, m.addr)
	if err != nil {
		return nil, err
	}
	if resp.Code != http.StatusOK {
		return nil, err
	}

	var parser expfmt.TextParser
	m2, err := parser.TextToMetricFamilies(bytes.NewBuffer(resp.Body))
	if err != nil {
		return nil, err
	}

	return &Data{Metric: m2}, nil
}
