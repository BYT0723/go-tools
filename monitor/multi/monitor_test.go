package multi

import (
	"context"
	"testing"
	"time"

	"github.com/BYT0723/go-tools/monitor"
	"github.com/stretchr/testify/assert"
)

type mockMonitor struct {
	ch      chan *monitor.Alert
	started bool
}

func newMockMonitor() *mockMonitor {
	return &mockMonitor{
		ch: make(chan *monitor.Alert, 10),
	}
}

func (m *mockMonitor) Start(ctx context.Context) {
	m.started = true
	go func() {
		<-ctx.Done()
		close(m.ch)
	}()
}

func (m *mockMonitor) Stop(ctx context.Context) {}
func (m *mockMonitor) SetCycle(d time.Duration)  {}
func (m *mockMonitor) SetTimeout(d time.Duration) {}
func (m *mockMonitor) Subscribe() <-chan *monitor.Alert {
	return m.ch
}

func (m *mockMonitor) sendAlert(a *monitor.Alert) {
	m.ch <- a
}

func TestNewMultiMonitor(t *testing.T) {
	t.Run("创建MultiMonitor", func(t *testing.T) {
		m1 := newMockMonitor()
		m2 := newMockMonitor()
		mm := NewMultiMonitor(m1, m2)
		assert.NotNil(t, mm)
		assert.Len(t, mm.monitors, 2)
	})

	t.Run("空monitor列表", func(t *testing.T) {
		mm := NewMultiMonitor()
		assert.NotNil(t, mm)
		assert.Empty(t, mm.monitors)
	})
}

func TestMultiMonitorStart(t *testing.T) {
	t.Run("Start启动所有子monitor", func(t *testing.T) {
		m1 := newMockMonitor()
		m2 := newMockMonitor()
		mm := NewMultiMonitor(m1, m2)

		ctx, cf := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cf()

		mm.Start(ctx)

		// 等待 monitor 启动
		time.Sleep(50 * time.Millisecond)
		assert.True(t, m1.started)
		assert.True(t, m2.started)
	})

	t.Run("子monitor告警转发到MultiMonitor", func(t *testing.T) {
		m1 := newMockMonitor()
		mm := NewMultiMonitor(m1)

		ctx, cf := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cf()

		mm.Start(ctx)

		// 发送告警消息
		alert := monitor.NewAlert(monitor.SeverityWarning, monitor.SourceAlertRule, "test", nil)
		m1.sendAlert(alert)

		select {
		case received := <-mm.Subscribe():
			assert.Equal(t, alert, received)
		case <-time.After(200 * time.Millisecond):
			t.Fatal("没有收到告警")
		}
	})
}

func TestMultiMonitorContextCancel(t *testing.T) {
	t.Run("context取消后Stop被调用", func(t *testing.T) {
		m1 := newMockMonitor()
		mm := NewMultiMonitor(m1)

		ctx, cf := context.WithCancel(context.Background())
		mm.Start(ctx)
		time.Sleep(50 * time.Millisecond)
		cf()

		// 等待 context 传播和 Stop 调用
		time.Sleep(100 * time.Millisecond)

		select {
		case _, ok := <-mm.Subscribe():
			assert.False(t, ok)
		default:
			// channel 应该已关闭
		}
	})
}
