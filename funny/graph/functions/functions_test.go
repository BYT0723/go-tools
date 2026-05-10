package functions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinaryFunctionType(t *testing.T) {
	t.Run("BinaryFunction 类型测试", func(t *testing.T) {
		var f BinaryFunction = func(x, y float64) bool { return x == y }
		assert.NotNil(t, f)
	})
}

func TestLove(t *testing.T) {
	t.Run("Love 函数测试", func(t *testing.T) {
		f := Love()
		assert.NotNil(t, f)

		// 爱心中心点应在函数范围内
		assert.True(t, f(0, 0))
		// 远离中心的点不应在范围内
		assert.False(t, f(10, 10))
	})
}

func TestCircularLove(t *testing.T) {
	t.Run("CircularLove 函数测试", func(t *testing.T) {
		f := CircularLove()
		assert.NotNil(t, f)
	})
}

func TestRoseLine(t *testing.T) {
	t.Run("RoseLine 函数测试", func(t *testing.T) {
		f := RoseLine(4, 1.0)
		assert.NotNil(t, f)
	})
}
