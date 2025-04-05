package ping

import (
	"fmt"
	"time"

	"github.com/BYT0723/go-tools/monitor"
	probing "github.com/prometheus-community/pro-bing"
)

// 0-100
func PktLossGt(f float64) monitor.AlertRule[probing.Statistics] {
	return func(s *probing.Statistics) (*monitor.Alert, bool) {
		if s.PacketLoss >= f {
			return monitor.NewAlert(
				monitor.SeverityError,
				monitor.SourceAlertRule,
				fmt.Sprintf("host %s packet loss: %.2f%%", s.Addr, s.PacketLoss),
				s,
			), true
		}
		return nil, false
	}
}

func MaxRttGt(rtt time.Duration) monitor.AlertRule[probing.Statistics] {
	return func(s *probing.Statistics) (*monitor.Alert, bool) {
		if s.MaxRtt >= rtt {
			return monitor.NewAlert(
				monitor.SeverityError,
				monitor.SourceAlertRule,
				fmt.Sprintf("host %s max rtt: %s", s.Addr, s.MaxRtt),
				s,
			), true
		}
		return nil, false
	}
}
