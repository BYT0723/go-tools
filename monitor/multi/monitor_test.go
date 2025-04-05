package multi

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/BYT0723/go-tools/monitor"
	"github.com/BYT0723/go-tools/monitor/ping"
	"github.com/BYT0723/go-tools/monitor/web"
	probing "github.com/prometheus-community/pro-bing"
)

func TestMultiMonitor(t *testing.T) {
	pm := ping.NewMonitor(
		"8.8.8.8",
		ping.WithCycle(10*time.Second),
		ping.WithAlertRules(func(s *probing.Statistics) (*monitor.Alert, bool) {
			fmt.Printf("s: %v", s)
			return nil, false
		}),
	)
	hm := web.NewMonitor(
		http.MethodGet,
		"https://baidu.com",
		web.WithCycle(10*time.Second),
		web.WithAlertRules(func(s *web.Statistics) (*monitor.Alert, bool) {
			fmt.Println(s.Code, s.Header)
			return nil, false
		}),
	)

	mm := NewMultiMonitor(pm, hm)

	ctx, cf := context.WithTimeout(context.Background(), time.Minute)
	defer cf()

	mm.Start(ctx)

	for ar := range mm.Subscribe() {
		fmt.Printf("ar: %v\n", ar)
	}
}
