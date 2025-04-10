package component

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/BYT0723/go-tools/monitor"
)

type (
	AlertComponent[T any] struct {
		mutex sync.Mutex
		rules []*alertRuleWrapper[T]
		ch    chan *monitor.Alert
	}
	alertRuleWrapper[T any] struct {
		counter uint32       // 累计触发次数
		last    atomic.Value // 上次告警时间
		alert   bool         // 是否已告警
		rule    monitor.AlertRule[T]
	}
)

var (
	cumulativeTimes = 3
	alertInterval   = 30 * time.Minute
)

func (m *AlertComponent[T]) AddAlertRule(ars ...monitor.AlertRule[T]) {
	m.mutex.Lock()
	for _, ar := range ars {
		m.rules = append(m.rules, &alertRuleWrapper[T]{rule: ar})
	}
	m.mutex.Unlock()
}

func (m *AlertComponent[T]) Evaluate(s *T) []*monitor.Alert {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	result := make([]*monitor.Alert, 0)
	for _, ar := range m.rules {
		if a, b := ar.rule(s); b {
			ar.counter++
			// 累计触发次数大于3次则告警
			if ar.counter >= uint32(cumulativeTimes) {
				if l, ok := ar.last.Load().(time.Time); !ar.alert ||
					(ok && time.Since(l) >= alertInterval) {
					// 初次告警 or 多次告警并超过告警间隔
					result = append(result, a)
					ar.alert = true
					ar.last.Store(time.Now())
				}
			}
		} else {
			ar.counter = 0
			ar.alert = false
		}
	}
	return result
}
