package osx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTermSize(t *testing.T) {
	t.Run("GetTermSize 测试", func(t *testing.T) {
		t.Run("调用GetTermSize不panic", func(t *testing.T) {
			assert.NotPanics(t, func() { GetTermSize() })
		})
	})
}

func TestCharmapDecode(t *testing.T) {
	t.Run("CharmapDecode 测试", func(t *testing.T) {
		t.Run("UTF-8 code page 65001 返回原始数据", func(t *testing.T) {
			input := []byte("hello")
			output, err := CharmapDecode(65001, input)
			assert.Nil(t, err)
			assert.Equal(t, "hello", string(output))
		})

		t.Run("未知 code page 返回错误", func(t *testing.T) {
			_, err := CharmapDecode(99999, []byte("hello"))
			assert.NotNil(t, err)
			assert.Equal(t, "unknown OEM Code Page", err.Error())
		})

		t.Run("GBK code page 936", func(t *testing.T) {
			_, err := CharmapDecode(936, []byte("hello"))
			assert.Nil(t, err)
		})

		t.Run("Shift-JIS code page 932", func(t *testing.T) {
			_, err := CharmapDecode(932, []byte("hello"))
			assert.Nil(t, err)
		})
	})
}
