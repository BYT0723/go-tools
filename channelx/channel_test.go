package channelx

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTryIn(t *testing.T) {
	t.Run("TryIn 测试", func(t *testing.T) {
		t.Run("发送到缓冲channel成功", func(t *testing.T) {
			ch := make(chan int, 1)
			assert.True(t, TryIn(ch, 42))
			assert.Equal(t, 42, <-ch)
		})

		t.Run("发送到满的channel失败", func(t *testing.T) {
			ch := make(chan int, 1)
			ch <- 1
			assert.False(t, TryIn(ch, 42))
		})

		t.Run("发送到无缓冲channel失败", func(t *testing.T) {
			ch := make(chan int)
			assert.False(t, TryIn(ch, 42))
		})
	})
}

func TestTryOut(t *testing.T) {
	t.Run("TryOut 测试", func(t *testing.T) {
		t.Run("从有数据的channel读取成功", func(t *testing.T) {
			ch := make(chan int, 1)
			ch <- 42
			v, ok := TryOut(ch)
			assert.True(t, ok)
			assert.Equal(t, 42, v)
		})

		t.Run("从空的channel读取失败", func(t *testing.T) {
			ch := make(chan int, 1)
			v, ok := TryOut(ch)
			assert.False(t, ok)
			assert.Equal(t, 0, v)
		})
	})
}

func TestIn(t *testing.T) {
	t.Run("In 测试", func(t *testing.T) {
		t.Run("正常发送成功", func(t *testing.T) {
			ch := make(chan int, 1)
			err := In(context.Background(), ch, 42)
			assert.Nil(t, err)
			assert.Equal(t, 42, <-ch)
		})

		t.Run("context取消时返回错误", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			ch := make(chan int)
			err := In(ctx, ch, 42)
			assert.Equal(t, context.Canceled, err)
		})
	})
}

func TestOut(t *testing.T) {
	t.Run("Out 测试", func(t *testing.T) {
		t.Run("正常读取成功", func(t *testing.T) {
			ch := make(chan int, 1)
			ch <- 42
			v, err := Out(context.Background(), ch)
			assert.Nil(t, err)
			assert.Equal(t, 42, v)
		})

		t.Run("context取消时返回错误", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			ch := make(chan int)
			_, err := Out(ctx, ch)
			assert.Equal(t, context.Canceled, err)
		})
	})
}

func TestInTimeout(t *testing.T) {
	t.Run("InTimeout 测试", func(t *testing.T) {
		t.Run("超时内发送成功", func(t *testing.T) {
			ch := make(chan int, 1)
			err := InTimeout(ch, 42, time.Second)
			assert.Nil(t, err)
			assert.Equal(t, 42, <-ch)
		})

		t.Run("超时发送失败", func(t *testing.T) {
			ch := make(chan int)
			err := InTimeout(ch, 42, time.Millisecond)
			assert.Equal(t, context.DeadlineExceeded, err)
		})
	})
}

func TestOutTimeout(t *testing.T) {
	t.Run("OutTimeout 测试", func(t *testing.T) {
		t.Run("超时内读取成功", func(t *testing.T) {
			ch := make(chan int, 1)
			ch <- 42
			v, err := OutTimeout(ch, time.Second)
			assert.Nil(t, err)
			assert.Equal(t, 42, v)
		})

		t.Run("超时读取失败", func(t *testing.T) {
			ch := make(chan int)
			_, err := OutTimeout(ch, time.Millisecond)
			assert.Equal(t, context.DeadlineExceeded, err)
		})
	})
}

func TestInDeadline(t *testing.T) {
	t.Run("InDeadline 测试", func(t *testing.T) {
		t.Run("截止时间前发送成功", func(t *testing.T) {
			ch := make(chan int, 1)
			err := InDeadline(ch, 42, time.Now().Add(time.Second))
			assert.Nil(t, err)
			assert.Equal(t, 42, <-ch)
		})

		t.Run("截止时间后发送失败", func(t *testing.T) {
			ch := make(chan int)
			err := InDeadline(ch, 42, time.Now().Add(-time.Second))
			assert.Equal(t, context.DeadlineExceeded, err)
		})
	})
}

func TestOutDeadline(t *testing.T) {
	t.Run("OutDeadline 测试", func(t *testing.T) {
		t.Run("截止时间前读取成功", func(t *testing.T) {
			ch := make(chan int, 1)
			ch <- 42
			v, err := OutDeadline(ch, time.Now().Add(time.Second))
			assert.Nil(t, err)
			assert.Equal(t, 42, v)
		})

		t.Run("截止时间后读取失败", func(t *testing.T) {
			ch := make(chan int)
			_, err := OutDeadline(ch, time.Now().Add(-time.Second))
			assert.Equal(t, context.DeadlineExceeded, err)
		})
	})
}
