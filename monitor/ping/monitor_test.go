package ping

import (
	"testing"
	"time"

	"github.com/BYT0723/go-tools/monitor"
	probing "github.com/prometheus-community/pro-bing"
	"github.com/stretchr/testify/assert"
)

func TestNewMonitor(t *testing.T) {
	t.Run("默认配置", func(t *testing.T) {
		m := NewMonitor("127.0.0.1")
		assert.NotNil(t, m)
		assert.Equal(t, "127.0.0.1", m.addr)
		assert.GreaterOrEqual(t, m.count, minCount)
	})

	t.Run("WithCount", func(t *testing.T) {
		m := NewMonitor("127.0.0.1", WithCount(5))
		assert.Equal(t, 5, m.count)
	})

	t.Run("WithTimeout - timeout受count限制", func(t *testing.T) {
		m := NewMonitor("127.0.0.1", WithTimeout(5*time.Second))
		// 默认count=10, timeout = max(5s, 10*1s) = 10s
		assert.Equal(t, 10*time.Second, m.Timeout)
	})

	t.Run("WithICMP", func(t *testing.T) {
		m := NewMonitor("127.0.0.1", WithICMP())
		assert.True(t, m.icmp)
	})

	t.Run("实现Monitor接口", func(t *testing.T) {
		m := NewMonitor("127.0.0.1")
		var _ monitor.Monitor = m
	})
}

func TestPktLossGt(t *testing.T) {
	t.Run("丢包率超过阈值触发告警", func(t *testing.T) {
		rule := PktLossGt(10)
		s := &probing.Statistics{
			Addr:       "8.8.8.8",
			PacketLoss: 20,
		}
		alert, matched := rule(s)
		assert.True(t, matched)
		assert.NotNil(t, alert)
		assert.Equal(t, monitor.SeverityError, alert.Severity)
	})

	t.Run("丢包率低于阈值不触发告警", func(t *testing.T) {
		rule := PktLossGt(50)
		s := &probing.Statistics{
			Addr:       "8.8.8.8",
			PacketLoss: 10,
		}
		alert, matched := rule(s)
		assert.False(t, matched)
		assert.Nil(t, alert)
	})

	t.Run("丢包率等于阈值触发告警", func(t *testing.T) {
		rule := PktLossGt(50)
		s := &probing.Statistics{
			Addr:       "8.8.8.8",
			PacketLoss: 50,
		}
		alert, matched := rule(s)
		assert.True(t, matched)
		assert.NotNil(t, alert)
	})
}

func TestMaxRttGt(t *testing.T) {
	t.Run("RTT超过阈值触发告警", func(t *testing.T) {
		rule := MaxRttGt(100 * time.Millisecond)
		s := &probing.Statistics{
			Addr:   "8.8.8.8",
			MaxRtt: 200 * time.Millisecond,
		}
		alert, matched := rule(s)
		assert.True(t, matched)
		assert.NotNil(t, alert)
	})

	t.Run("RTT低于阈值不触发告警", func(t *testing.T) {
		rule := MaxRttGt(time.Second)
		s := &probing.Statistics{
			Addr:   "8.8.8.8",
			MaxRtt: 50 * time.Millisecond,
		}
		alert, matched := rule(s)
		assert.False(t, matched)
		assert.Nil(t, alert)
	})

	t.Run("RTT等于阈值触发告警", func(t *testing.T) {
		rule := MaxRttGt(100 * time.Millisecond)
		s := &probing.Statistics{
			Addr:   "8.8.8.8",
			MaxRtt: 100 * time.Millisecond,
		}
		alert, matched := rule(s)
		assert.True(t, matched)
		assert.NotNil(t, alert)
	})
}

func TestPingAlertComponent(t *testing.T) {
	t.Run("累积触发告警", func(t *testing.T) {
		m := NewMonitor("127.0.0.1")
		m.AddAlertRule(PktLossGt(10))

		s := &probing.Statistics{
			Addr:       "127.0.0.1",
			PacketLoss: 50,
		}

		// 第1-2次匹配但累积次数不够，不触发
		for i := 0; i < 2; i++ {
			alerts := m.Evaluate(s)
			assert.Empty(t, alerts)
		}

		// 第3次触发告警
		alerts := m.Evaluate(s)
		assert.NotEmpty(t, alerts)
	})
}
