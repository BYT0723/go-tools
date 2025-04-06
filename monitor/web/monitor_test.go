package web

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/BYT0723/go-tools/monitor"
)

func TestHttpMonitor(t *testing.T) {
	m := NewMonitor(http.MethodGet, "https://www.baidu.com", WithTimeout(5*time.Second))
	m.SetCycle(10 * time.Second)
	m.AddAlertRule(func(s *Statistics) (*monitor.Alert, bool) {
		fmt.Println(s.Code, s.Header)
		if s.Code != http.StatusOK {
			return monitor.NewAlert(
				monitor.SeverityWarning,
				monitor.SourceAlertRule,
				"status code not 200",
				s,
			), true
		}
		return nil, false
	})

	ctx, cf := context.WithTimeout(context.Background(), time.Minute)
	defer cf()

	m.Start(ctx)

	for ar := range m.Subscribe() {
		fmt.Printf("ar: %v\n", ar)
	}
}
