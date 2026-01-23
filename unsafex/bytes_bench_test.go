package unsafex

import (
	"testing"
)

var (
	testString   = "Hello, World! This is a test string for benchmarking."
	testBytes    = []byte("Hello, World! This is a test string for benchmarking.")
	resultBytes  []byte
	resultString string
)

func BenchmarkUnsafeBytes(b *testing.B) {
	var r []byte
	for i := 0; i < b.N; i++ {
		r = UnsafeBytes(testString)
	}
	resultBytes = r
}

func BenchmarkStandardBytes(b *testing.B) {
	var r []byte
	for i := 0; i < b.N; i++ {
		r = []byte(testString)
	}
	resultBytes = r
}

func BenchmarkUnsafeString(b *testing.B) {
	var r string
	for i := 0; i < b.N; i++ {
		r = UnsafeString(testBytes)
	}
	resultString = r
}

func BenchmarkStandardString(b *testing.B) {
	var r string
	for i := 0; i < b.N; i++ {
		r = string(testBytes)
	}
	resultString = r
}

// Benchmark for small strings/bytes
var (
	smallString = "tiny"
	smallBytes  = []byte("tiny")
)

func BenchmarkUnsafeBytesSmall(b *testing.B) {
	var r []byte
	for i := 0; i < b.N; i++ {
		r = UnsafeBytes(smallString)
	}
	resultBytes = r
}

func BenchmarkStandardBytesSmall(b *testing.B) {
	var r []byte
	for i := 0; i < b.N; i++ {
		r = []byte(smallString)
	}
	resultBytes = r
}

func BenchmarkUnsafeStringSmall(b *testing.B) {
	var r string
	for i := 0; i < b.N; i++ {
		r = UnsafeString(smallBytes)
	}
	resultString = r
}

func BenchmarkStandardStringSmall(b *testing.B) {
	var r string
	for i := 0; i < b.N; i++ {
		r = string(smallBytes)
	}
	resultString = r
}

// Benchmark for empty strings/bytes
var (
	emptyString = ""
	emptyBytes  = []byte{}
)

func BenchmarkUnsafeBytesEmpty(b *testing.B) {
	var r []byte
	for i := 0; i < b.N; i++ {
		r = UnsafeBytes(emptyString)
	}
	resultBytes = r
}

func BenchmarkStandardBytesEmpty(b *testing.B) {
	var r []byte
	for i := 0; i < b.N; i++ {
		r = []byte(emptyString)
	}
	resultBytes = r
}

func BenchmarkUnsafeStringEmpty(b *testing.B) {
	var r string
	for i := 0; i < b.N; i++ {
		r = UnsafeString(emptyBytes)
	}
	resultString = r
}

func BenchmarkStandardStringEmpty(b *testing.B) {
	var r string
	for i := 0; i < b.N; i++ {
		r = string(emptyBytes)
	}
	resultString = r
}
