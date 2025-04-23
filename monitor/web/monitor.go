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
	Monitor struct {
		*component.MonitorComponent
		component.AlertComponent[httpx.Response]

		client *httpx.Client

		method  string
		addr    string
		header  http.Header
		payload any
	}
	Option func(*Monitor)
)

var (
	defaultClient = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
	defautlCycle = time.Minute
)

func NewMonitor(method string, addr string, opts ...Option) *Monitor {
	p := &Monitor{
		MonitorComponent: component.NewMonitorComponent(),
		client:           httpx.NewClient(httpx.WithHttpClient(&defaultClient)),
		addr:             addr,
	}
	p.SetCycle(defautlCycle)
	p.Timeout = defautlCycle

	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (m *Monitor) Start(ctx context.Context) {
	m.SetContext(ctx)

	go func() {
		t := m.MonitorComponent.Ticker()
		defer t.Stop()

		s, err := m.do()
		if err != nil {
			m.Notify(monitor.InternalAlert(err))
		} else {
			m.Notify(m.Evaluate(s)...)
		}

		for {
			select {
			case <-m.Context().Done():
				m.MonitorComponent.Stop(ctx)
				return
			case <-t.C:
				s, err := m.do()
				if err != nil {
					m.Notify(monitor.InternalAlert(err))
					continue
				}
				m.Notify(m.Evaluate(s)...)
			}
		}
	}()
}

func (m *Monitor) do() (*httpx.Response, error) {
	ctx, cf := context.WithTimeout(context.Background(), m.Timeout)
	defer cf()
	return m.client.Do(
		ctx,
		m.method,
		m.addr,
		httpx.WithHeader(m.header),
		httpx.WithPayload(m.payload),
	)
}
