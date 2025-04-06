package ping

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/BYT0723/go-tools/monitor"
	probing "github.com/prometheus-community/pro-bing"
)

func TestPingMonitor(t *testing.T) {
	m := NewMonitor("8.8.8.8", WithCount(5))
	m.SetCycle(10 * time.Second)
	m.AddAlertRule(func(s *probing.Statistics) (*monitor.Alert, bool) {
		fmt.Printf("s: %v\n", s)
		if s.PacketLoss > 0 {
			return monitor.NewAlert(
				monitor.SeverityWarning,
				monitor.SourceAlertRule,
				fmt.Sprintf("host %s packet loss: %.2f%%", s.Addr, s.PacketLoss),
				s,
			), true
		}
		return nil, false
	})

	ctx, cf := context.WithTimeout(context.Background(), time.Minute)
	defer cf()

	m.Start(ctx)

	for ar := range m.Subscribe() {
		fmt.Printf("ar: %+v\n", ar)
	}
}
