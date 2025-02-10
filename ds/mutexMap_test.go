package ds

import (
	"fmt"
	"testing"
)

type Item struct {
	value string
}

func TestMutexMap(t *testing.T) {
	m := NewMutexMap[int, Item]()

	i1 := Item{"abc"}
	i2 := Item{"def"}

	m.Store(1, i1)
	m.Store(2, i2)

	fmt.Printf("m.entries: %v\n", m.entries)
}
