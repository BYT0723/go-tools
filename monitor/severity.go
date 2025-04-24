package monitor

type Severity uint8

const (
	SeverityDebug    Severity = iota // Debugging
	SeverityInfo                     // Information
	SeverityWarning                  // Warning
	SeverityError                    // Error
	SeverityCritical                 // Critical
)

// String implements fmt.String interface
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
