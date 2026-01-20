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
}
