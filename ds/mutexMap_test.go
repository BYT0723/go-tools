package ds

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type Item struct {
	value string
}

func TestMutexMap(t *testing.T) {
	var (
		m   = NewMutexMap[int, Item]()
		wg  sync.WaitGroup
		i1  = Item{"abc"}
		i2  = Item{"def"}
		ch1 = make(chan struct{}, 1)
		ch2 = make(chan struct{}, 1)
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		m.Store(1, i1)
		ch1 <- struct{}{}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		m.Store(2, i2)
		ch2 <- struct{}{}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ch2
		m.Store(2, i1)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ch1
		m.Store(1, i2)
	}()

	wg.Wait()

	require.Equal(t, 2, m.Len())
	require.Equal(t, []Item{i2, i1}, m.Values())
}
