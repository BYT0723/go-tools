package component

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/BYT0723/go-tools/monitor"
)

type (
	// AlertComponent is a generic component that handles alerting logic by evaluating rules.
	// It stores a list of alert rules and triggers alerts based on those rules.
	AlertComponent[T any] struct {
		mutex sync.Mutex             // Mutex to protect concurrent access to rules
		rules []*alertRuleWrapper[T] // A list of alert rule wrappers
	}

	// alertRuleWrapper wraps an alert rule with additional state for tracking
	// the number of times a rule was triggered and the last alert time.
	alertRuleWrapper[T any] struct {
		counter uint32               // The number of times this rule has been triggered
		last    atomic.Value         // The last time this alert was triggered
		alert   bool                 // Indicates whether an alert has been triggered
		rule    monitor.AlertRule[T] // The alert rule function
	}
)

var (
	// cumulativeTimes is the threshold for the number of times a rule must be triggered
	// before an alert is sent.
	cumulativeTimes = 3

	// alertInterval defines the time window in which the rule must not have triggered an alert
	// to allow another alert to be sent.
	alertInterval = 30 * time.Minute
)

// AddAlertRule adds one or more alert rules to the AlertComponent.
func (m *AlertComponent[T]) AddAlertRule(ars ...monitor.AlertRule[T]) {
	m.mutex.Lock()
	// Add each rule to the component's list of rules
	for _, ar := range ars {
		m.rules = append(m.rules, &alertRuleWrapper[T]{rule: ar})
	}
	m.mutex.Unlock()
}

// Evaluate evaluates the alert rules against the provided data (s).
// If any rule triggers, an alert is created and returned.
// It returns a slice of alerts that were triggered during the evaluation.
func (m *AlertComponent[T]) Evaluate(s *T) []*monitor.Alert {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Initialize a slice to store triggered alerts
	result := make([]*monitor.Alert, 0)

	// Iterate over each rule and check if it triggers
	for _, ar := range m.rules {
		// Evaluate the rule
		a, matched := ar.rule(s)
		if !matched {
			continue
		}
		if a != nil {
			// Increment the trigger counter
			ar.counter++
			// If the rule has been triggered enough times, check the alert conditions
			if ar.counter >= uint32(cumulativeTimes) {
				// Check if the last alert time has passed the interval or if it's the first alert
				if l, ok := ar.last.Load().(time.Time); !ar.alert ||
					(ok && time.Since(l) >= alertInterval) {
					// Trigger the alert if the conditions are met
					result = append(result, a)
					ar.alert = true
					ar.last.Store(time.Now()) // Store the time of this alert
				}
			}
		} else {
			// Reset counter and alert state if the rule no longer triggers
			ar.counter = 0
			ar.alert = false
		}
	}
	// Return the list of triggered alerts
	return result
}
