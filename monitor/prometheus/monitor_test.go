package prometheus

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/BYT0723/go-tools/monitor"
	"github.com/stretchr/testify/assert"
)

func TestNewMonitor(t *testing.T) {
	t.Run("默认配置", func(t *testing.T) {
		m := NewMonitor("http://localhost:9100/metrics")
		assert.NotNil(t, m)
		assert.Equal(t, "http://localhost:9100/metrics", m.addr)
	})

	t.Run("WithCycle", func(t *testing.T) {
		m := NewMonitor("http://localhost:9100/metrics", WithCycle(30*time.Second))
		assert.NotNil(t, m)
	})

	t.Run("WithTimeout", func(t *testing.T) {
		m := NewMonitor("http://localhost:9100/metrics", WithTimeout(3*time.Second))
		assert.Equal(t, 3*time.Second, m.Timeout)
	})

	t.Run("实现Monitor接口", func(t *testing.T) {
		m := NewMonitor("http://localhost:9100/metrics")
		var _ monitor.Monitor = m
	})
}

func TestDataStruct(t *testing.T) {
	t.Run("Data 结构创建", func(t *testing.T) {
		d := &Data{Metric: make(map[string]*MetricFamily)}
		assert.NotNil(t, d)
		assert.NotNil(t, d.Metric)
	})
}

func TestPrometheusAlertRule(t *testing.T) {
	t.Run("告警规则可以正常添加和评估", func(t *testing.T) {
		m := NewMonitor("http://localhost:9100/metrics")

		m.AddAlertRule(func(d *Data) (*monitor.Alert, bool) {
			if d.Metric != nil {
				return monitor.NewAlert(monitor.SeverityInfo, monitor.SourceAlertRule, "has metrics", nil), true
			}
			return nil, false
		})

		d := &Data{Metric: make(map[string]*MetricFamily)}
		alerts := m.Evaluate(d)

		// 首次匹配但累积次数不够，不触发告警
		assert.Empty(t, alerts)
	})
}

func TestPrometheusWithServer(t *testing.T) {
	t.Run("Start/Stop 基本流程", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(""))
		}))
		defer srv.Close()

		m := NewMonitor(srv.URL, WithTimeout(time.Second))
		m.SetCycle(50 * time.Millisecond)

		ctx, cf := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer cf()

		go m.Start(ctx)

		<-ctx.Done()
		time.Sleep(50 * time.Millisecond)
	})
}
