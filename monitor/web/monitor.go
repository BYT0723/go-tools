package web

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/BYT0723/go-tools/monitor"
	"github.com/BYT0723/go-tools/transport/httpx"
)

type (
	Monitor struct {
		client  *httpx.Client
		decoder httpx.Decoder

		method  string
		addr    string
		header  http.Header
		payload any

		cycle   time.Duration
		timeout time.Duration
		ctx     context.Context
		cf      context.CancelFunc

		alertRulesMutex sync.Mutex
		alertRules      []monitor.AlertRule[Statistics]
		ch              chan *monitor.Alert
	}
	Statistics struct {
		Code    int
		Header  http.Header
		Resp    []byte
		Payload map[string]any
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
	cycle   = time.Minute
	timeout = 5 * time.Second
)

func NewMonitor(method string, addr string, opts ...Option) *Monitor {
	p := &Monitor{
		client:  httpx.NewClient(httpx.WithHttpClient(&defaultClient)),
		addr:    addr,
		cycle:   cycle,
		timeout: timeout,
		ch:      make(chan *monitor.Alert, 1024),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (m *Monitor) Start(ctx context.Context) {
	m.ctx, m.cf = context.WithCancel(ctx)

	go func() {
		ticker := time.NewTicker(m.cycle)
		defer ticker.Stop()

		s, err := m.do()
		if err != nil {
			m.ch <- monitor.InternalAlert(err)
		} else {
			m.alertRulesMutex.Lock()
			for _, ar := range m.alertRules {
				if a, b := ar(s); b {
					m.ch <- a
				}
			}
			m.alertRulesMutex.Unlock()
		}

		for {
			select {
			case <-m.ctx.Done():
				m.cf()
				return
			case <-ticker.C:
				s, err := m.do()
				if err != nil {
					m.ch <- monitor.InternalAlert(err)
					continue
				}
				m.alertRulesMutex.Lock()
				for _, ar := range m.alertRules {
					if a, b := ar(s); b {
						m.ch <- a
					}
				}
				m.alertRulesMutex.Unlock()
			}
		}
	}()
}

func (m *Monitor) Stop() {
	if m.cf != nil {
		m.cf()
	}
	close(m.ch)
}

func (m *Monitor) Subscribe() <-chan *monitor.Alert {
	return m.ch
}

func (m *Monitor) AddAlertRule(ars ...monitor.AlertRule[Statistics]) {
	m.alertRulesMutex.Lock()
	m.alertRules = append(m.alertRules, ars...)
	m.alertRulesMutex.Unlock()
}

func (m *Monitor) do() (*Statistics, error) {
	ctx, cf := context.WithTimeout(m.ctx, m.timeout)
	defer cf()
	code, header, rc, err := m.client.Do(ctx, m.method, m.addr, m.header, m.payload)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	b, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	stats := &Statistics{Code: code, Header: header, Resp: b}

	if m.decoder != nil {
		payload := make(map[string]any)
		if err := m.decoder(ctx, rc, &payload); err == nil {
			stats.Payload = payload
		}
	}

	return stats, nil
}
