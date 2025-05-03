package prometheus

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/BYT0723/go-tools/monitor"
)

func TestPrometheusMonitor(t *testing.T) {
	m := NewMonitor(
		"http://localhost:9100/metrics",
		WithCycle(20*time.Second),
		WithTimeout(5*time.Second),
	)

	m.AddAlertRule(func(d *Data) (alert *monitor.Alert, matched bool) {
		if mf, ok := d.Metric["go_memstats_frees_total"]; ok {
			matched = true
			for _, m := range mf.GetMetric() {
				if c := m.GetCounter(); c != nil {
					if v := c.GetValue(); v > 0 {
						alert = monitor.NewAlert(monitor.SeverityCritical, monitor.SourceAlertRule, "go_memstats_frees_total > 0", v)
						break
					}
				}
			}
		}
		return
	})

	ctx, cf := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cf()

	m.Start(ctx)

	for ar := range m.Subscribe() {
		fmt.Printf("ar: %v\n", ar)
	}
}
