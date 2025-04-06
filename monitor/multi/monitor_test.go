package multi

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/BYT0723/go-tools/monitor/ping"
	"github.com/BYT0723/go-tools/monitor/web"
)

func TestMultiMonitor(t *testing.T) {
	pm := ping.NewMonitor("8.8.8.8")
	pm.SetCycle(10 * time.Second)
	pm.AddAlertRule(
		ping.MaxRttGt(10*time.Millisecond),
		ping.PktLossGt(10),
	)

	hm := web.NewMonitor(http.MethodGet, "https://baidu.com")
	hm.SetCycle(10 * time.Second)
	hm.AddAlertRule(web.CodeEqual(200))

	mm := NewMultiMonitor(pm, hm)

	ctx, cf := context.WithTimeout(context.Background(), time.Minute)
	defer cf()

	mm.Start(ctx)

	for ar := range mm.Subscribe() {
		log.Printf("ar: %v\n", ar)
	}
}
