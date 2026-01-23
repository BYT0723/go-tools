package unsafex

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestUnsafeBytes(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var str string
		b := UnsafeBytes(str)

		assert.Equal(t, unsafe.StringData(str), unsafe.SliceData(b))
		assert.Equal(t, len(str), len(b))
		assert.Equal(t, str, string(b))
	})

	t.Run("string", func(t *testing.T) {
		str := "hello"
		b := UnsafeBytes(str)

		assert.Equal(t, unsafe.StringData(str), unsafe.SliceData(b))
		assert.Equal(t, len(str), len(b))
		assert.Equal(t, str, string(b))
	})

	t.Run("unicode", func(t *testing.T) {
		str := "Hello, 世界!"
		b := UnsafeBytes(str)

		assert.Equal(t, unsafe.StringData(str), unsafe.SliceData(b))
		assert.Equal(t, len(str), len(b))
		assert.Equal(t, str, string(b))
	})

	t.Run("zero length non-nil", func(t *testing.T) {
		str := ""
		b := UnsafeBytes(str)

		assert.Equal(t, unsafe.StringData(str), unsafe.SliceData(b))
		assert.Equal(t, len(str), len(b))
		assert.Equal(t, str, string(b))
	})
}

func TestUnsafeString(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var b []byte
		str := UnsafeString(b)

		assert.Equal(t, unsafe.StringData(str), unsafe.SliceData(b))
		assert.Equal(t, len(str), len(b))
		assert.Equal(t, str, string(b))
	})

	t.Run("bytes", func(t *testing.T) {
		var b []byte = []byte{'h', 'e', 'l', 'l', 'o'}
		str := UnsafeString(b)

		assert.Equal(t, unsafe.StringData(str), unsafe.SliceData(b))
		assert.Equal(t, len(str), len(b))
		assert.Equal(t, str, string(b))
	})

	t.Run("unicode bytes", func(t *testing.T) {
		b := []byte("Hello, 世界!")
		str := UnsafeString(b)

		assert.Equal(t, unsafe.StringData(str), unsafe.SliceData(b))
		assert.Equal(t, len(str), len(b))
		assert.Equal(t, str, string(b))
	})

	t.Run("zero length non-nil", func(t *testing.T) {
		b := []byte{}
		str := UnsafeString(b)

		assert.Equal(t, unsafe.StringData(str), unsafe.SliceData(b))
		assert.Equal(t, len(str), len(b))
		assert.Equal(t, str, string(b))
	})

	t.Run("capacity greater than length", func(t *testing.T) {
		b := make([]byte, 3, 10)
		b[0] = 'a'
		b[1] = 'b'
		b[2] = 'c'
		str := UnsafeString(b)

		// String length should be 3, not 10
		assert.Equal(t, 3, len(str))
		assert.Equal(t, "abc", str)
		assert.Equal(t, unsafe.StringData(str), unsafe.SliceData(b))
	})
}
