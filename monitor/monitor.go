package monitor

import (
	"context"
	"time"
)

type (
	Monitor interface {
		Start(context.Context)
		Stop(context.Context)
		Subscribe() <-chan *Alert
	}
	Alert struct {
		Ts       time.Time
		Severity Severity
		Source   Source
		Err      error
		Descr    string
		Payload  any
	}
	AlertRule[T any] func(*T) (*Alert, bool)
)

func InternalAlert(err error) *Alert {
	return &Alert{Ts: time.Now(), Severity: SeverityError, Source: SourceInternal, Err: err}
}

func NewAlert(severity Severity, source Source, descr string, payload any) *Alert {
	return &Alert{
		Ts:       time.Now(),
		Severity: severity,
		Source:   source,
		Descr:    descr,
		Payload:  payload,
	}
}
