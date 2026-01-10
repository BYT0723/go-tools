package ds

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFastBus(t *testing.T) {
	var (
		b  = NewFastHub[string](3)
		wg sync.WaitGroup
	)

	for range 10 {
		s, _ := b.Subscribe()
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer b.Unsubscribe(s)
			v := <-s.Channel()
			assert.Equal(t, "string1", v)
			v = <-s.Channel()
			assert.Equal(t, "string2", v)
			v = <-s.Channel()
			assert.Equal(t, "", v)
		}()
	}

	b.Publish("string1")
	b.Publish("string2")
	b.Publish("")

	b.Close()

	wg.Wait()
}
