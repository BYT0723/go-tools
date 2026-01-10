package ds

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderBus(t *testing.T) {
	var (
		b  = NewOrderHub[string](10)
		wg sync.WaitGroup
	)

	for range 10 {
		s, _ := b.Subscribe()
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer b.Unsubscribe(s)
			v := <-s.Channel()
			assert.Equal(t, "hello", v)
			v = <-s.Channel()
			assert.Equal(t, "my name is walter", v)
		}()
	}

	b.Publish("hello")
	b.Publish("my name is walter")

	b.Close()

	wg.Wait()
}
