package set

import (
	"fmt"
	"math/rand"
	"testing"
)

type Payload struct {
	Field1 int
	Field2 string
}

var (
	s1         = NewSet[*Payload]()
	s2         = NewSet[*Payload]()
	randString = func() string {
		var (
			n        = 10 + rand.Intn(10)
			strBytes = make([]byte, n)
		)
		for i := range strBytes {
			strBytes[i] = byte(65 + rand.Intn(26))
		}
		return string(strBytes)
	}
)

func BenchmarkAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s1.Append(&Payload{
			Field1: rand.Intn(10),
			Field2: randString(),
		})
	}
}

func BenchmarkAppend2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s2.Append2(&Payload{
			Field1: rand.Intn(10),
			Field2: randString(),
		})
	}
}

func BenchmarkLength(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s1.Length()
	}
}

func BenchmarkLength2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s2.Length2()
	}
}

func BenchmarkUnion(b *testing.B) {
	fmt.Printf("s1.Length(): %v\n", s1.Length())
	fmt.Printf("s2.Length2(): %v\n", s2.Length2())
	for i := 0; i < b.N; i++ {
		s1.Union(s2)
	}
}
