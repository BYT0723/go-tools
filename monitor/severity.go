package monitor

type Severity uint8

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
