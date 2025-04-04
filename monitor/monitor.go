package monitor

import "time"

type (
	Severity uint8
	Source   string
	Alert    struct {
		Ts       time.Time
		Severity Severity
		Source   Source
		Err      error
		Descr    string
		Payload  any
	}
	AlertRule[T any] func(*T) (*Alert, bool)
)

const (
	SeverityDebug    Severity = iota // 调试信息
	SeverityInfo                     // 一般信息
	SeverityWarning                  // 警告
	SeverityError                    // 错误
	SeverityCritical                 // 严重错误
)

// String 实现Severity的字符串表示
func (s Severity) String() string {
	switch s {
	case SeverityDebug:
		return "DEBUG"
	case SeverityInfo:
		return "INFO"
	case SeverityWarning:
		return "WARNING"
	case SeverityError:
		return "ERROR"
	case SeverityCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

const (
	SourceInternal  Source = "internal"
	SourceAlertRule Source = "alert_rule"
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
