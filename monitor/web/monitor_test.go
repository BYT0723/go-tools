package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/BYT0723/go-tools/monitor"
	"github.com/BYT0723/go-tools/transport/httpx"
	"github.com/stretchr/testify/assert"
)

func TestNewMonitor(t *testing.T) {
	t.Run("默认配置", func(t *testing.T) {
		m := NewMonitor(http.MethodGet, "http://example.com")
		assert.NotNil(t, m)
		assert.Equal(t, "http://example.com", m.addr)
	})

	t.Run("WithTimeout", func(t *testing.T) {
		m := NewMonitor(http.MethodGet, "http://example.com", WithTimeout(3*time.Second))
		assert.Equal(t, 3*time.Second, m.Timeout)
	})

	t.Run("实现Monitor接口", func(t *testing.T) {
		m := NewMonitor(http.MethodGet, "http://example.com")
		var _ monitor.Monitor = m
	})
}

func TestCodeEqual(t *testing.T) {
	t.Run("状态码匹配触发告警", func(t *testing.T) {
		rule := CodeEqual(500)
		resp := &httpx.Response{Code: 500}
		alert, matched := rule(resp)
		assert.True(t, matched)
		assert.NotNil(t, alert)
	})

	t.Run("状态码不匹配不触发告警", func(t *testing.T) {
		rule := CodeEqual(500)
		resp := &httpx.Response{Code: 200}
		alert, matched := rule(resp)
		assert.False(t, matched)
		assert.Nil(t, alert)
	})
}

func TestCodeNotEqual(t *testing.T) {
	t.Run("状态码不等于目标值触发告警", func(t *testing.T) {
		rule := CodeNotEqual(200)
		resp := &httpx.Response{Code: 500}
		alert, matched := rule(resp)
		assert.True(t, matched)
		assert.NotNil(t, alert)
	})

	t.Run("状态码等于目标值不触发告警", func(t *testing.T) {
		rule := CodeNotEqual(200)
		resp := &httpx.Response{Code: 200}
		alert, matched := rule(resp)
		assert.False(t, matched)
		assert.Nil(t, alert)
	})
}

func TestHeaderContains(t *testing.T) {
	t.Run("Header包含指定值时触发告警", func(t *testing.T) {
		rule := HeaderContains("X-Custom", "error")
		resp := &httpx.Response{
			Header: http.Header{"X-Custom": {"error"}},
		}
		alert, matched := rule(resp)
		assert.True(t, matched)
		assert.NotNil(t, alert)
	})

	t.Run("Header不包含指定值时不触发告警", func(t *testing.T) {
		rule := HeaderContains("X-Custom", "error")
		resp := &httpx.Response{
			Header: http.Header{"X-Custom": {"ok"}},
		}
		alert, matched := rule(resp)
		assert.False(t, matched)
		assert.Nil(t, alert)
	})
}

func TestHeaderNotContains(t *testing.T) {
	t.Run("Header不包含指定值时触发告警", func(t *testing.T) {
		rule := HeaderNotContains("Server", "nginx")
		resp := &httpx.Response{
			Header: http.Header{"Server": {"apache"}},
		}
		alert, matched := rule(resp)
		assert.True(t, matched)
		assert.NotNil(t, alert)
	})

	t.Run("Header包含指定值时不触发告警", func(t *testing.T) {
		rule := HeaderNotContains("Server", "nginx")
		resp := &httpx.Response{
			Header: http.Header{"Server": {"nginx"}},
		}
		alert, matched := rule(resp)
		assert.False(t, matched)
		assert.Nil(t, alert)
	})
}

func TestWebAlertComponent(t *testing.T) {
	t.Run("累积触发告警", func(t *testing.T) {
		m := NewMonitor(http.MethodGet, "http://example.com")
		m.AddAlertRule(CodeNotEqual(200))

		resp := &httpx.Response{Code: 500}

		for i := 0; i < 2; i++ {
			alerts := m.Evaluate(resp)
			assert.Empty(t, alerts)
		}

		alerts := m.Evaluate(resp)
		assert.NotEmpty(t, alerts)
	})
}

func TestWebMonitorWithServer(t *testing.T) {
	t.Run("Start/Stop 基本流程", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer srv.Close()

		m := NewMonitor(http.MethodGet, srv.URL, WithTimeout(time.Second))
		m.SetCycle(50 * time.Millisecond)

		ctx, cf := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer cf()

		go m.Start(ctx)

		// 等待 context 超时自动停止
		<-ctx.Done()
		time.Sleep(50 * time.Millisecond)
	})
}
