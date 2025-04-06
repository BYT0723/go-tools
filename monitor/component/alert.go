package component

import (
	"sync"

	"github.com/BYT0723/go-tools/monitor"
)

type (
	AlertComponent[T any] struct {
		mutex sync.Mutex
		rules []*alertRuleWrapper[T]
		ch    chan *monitor.Alert
	}
	alertRuleWrapper[T any] struct {
		counter uint32
		rule    monitor.AlertRule[T]
	}
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
			if ar.counter >= 3 {
				result = append(result, a)
			}
		} else {
			ar.counter = 0
		}
	}
	return result
}
