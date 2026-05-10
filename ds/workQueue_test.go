package ds

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkQueueAddTopic(t *testing.T) {
	t.Run("WorkQueue AddTopic 测试", func(t *testing.T) {
		wq := NewWorkQueue[int]()

		t.Run("添加新topic成功", func(t *testing.T) {
			err := wq.AddTopic("topic1", 10)
			assert.Nil(t, err)
		})

		t.Run("重复添加同一topic返回错误", func(t *testing.T) {
			wq.AddTopic("topic1", 10)
			err := wq.AddTopic("topic1", 10)
			assert.Equal(t, ErrTopicAlreadyExists, err)
		})
	})
}

func TestWorkQueueRemoveTopic(t *testing.T) {
	t.Run("WorkQueue RemoveTopic 测试", func(t *testing.T) {
		t.Run("移除已存在的topic", func(t *testing.T) {
			wq := NewWorkQueue[int]()
			wq.AddTopic("topic1", 10)
			err := wq.RemoveTopic("topic1")
			assert.Nil(t, err)
		})

		t.Run("移除不存在的topic不返回错误", func(t *testing.T) {
			wq := NewWorkQueue[int]()
			err := wq.RemoveTopic("nonexistent")
			assert.Nil(t, err)
		})
	})
}

func TestWorkQueuePublish(t *testing.T) {
	t.Run("WorkQueue Publish 测试", func(t *testing.T) {
		t.Run("发布消息到已有topic", func(t *testing.T) {
			wq := NewWorkQueue[int]()
			wq.AddTopic("topic1", 10)
			err := wq.Publish("topic1", 42)
			assert.Nil(t, err)
		})

		t.Run("发布消息到不存在topic返回错误", func(t *testing.T) {
			wq := NewWorkQueue[int]()
			err := wq.Publish("nonexistent", 42)
			assert.Equal(t, ErrTopicNotFound, err)
		})

		t.Run("发布到满队列返回队列满错误", func(t *testing.T) {
			wq := NewWorkQueue[int]()
			wq.AddTopic("topic1", 1)
			wq.Publish("topic1", 1)
			err := wq.Publish("topic1", 2)
			assert.Equal(t, ErrTopicQueueFull, err)
		})
	})
}

func TestWorkQueueSubscribe(t *testing.T) {
	t.Run("WorkQueue Subscribe 测试", func(t *testing.T) {
		t.Run("订阅已存在的topic成功", func(t *testing.T) {
			wq := NewWorkQueue[int]()
			wq.AddTopic("topic1", 10)
			ch, err := wq.Subscribe("topic1")
			assert.Nil(t, err)
			assert.NotNil(t, ch)
		})

		t.Run("订阅不存在的topic返回错误", func(t *testing.T) {
			wq := NewWorkQueue[int]()
			ch, err := wq.Subscribe("nonexistent")
			assert.Equal(t, ErrTopicNotFound, err)
			assert.Nil(t, ch)
		})
	})
}

func TestWorkQueuePubSub(t *testing.T) {
	t.Run("WorkQueue 发布订阅消息传递测试", func(t *testing.T) {
		t.Run("订阅者能收到发布的消息", func(t *testing.T) {
			wq := NewWorkQueue[int]()
			wq.AddTopic("topic1", 10)

			ch, _ := wq.Subscribe("topic1")
			var wg sync.WaitGroup
			wg.Add(1)
			var result int
			go func() {
				defer wg.Done()
				result = <-ch
			}()

			wq.Publish("topic1", 42)
			wg.Wait()
			assert.Equal(t, 42, result)
		})

		t.Run("两个topic互不影响", func(t *testing.T) {
			wq := NewWorkQueue[int]()
			wq.AddTopic("topic1", 10)
			wq.AddTopic("topic2", 10)

			ch1, _ := wq.Subscribe("topic1")
			ch2, _ := wq.Subscribe("topic2")

			wq.Publish("topic1", 1)
			wq.Publish("topic2", 2)

			assert.Equal(t, 1, <-ch1)
			assert.Equal(t, 2, <-ch2)
		})
	})
}

func TestWorkQueueClose(t *testing.T) {
	t.Run("WorkQueue Close 测试", func(t *testing.T) {
		t.Run("Close 后channel被关闭", func(t *testing.T) {
			wq := NewWorkQueue[int]()
			wq.AddTopic("topic1", 10)
			wq.Close()

			ch, _ := wq.Subscribe("topic1")
			_, ok := <-ch
			assert.False(t, ok)
		})
	})
}

func TestWorkQueueErrors(t *testing.T) {
	t.Run("WorkQueue 错误变量测试", func(t *testing.T) {
		assert.Equal(t, "topic not found", ErrTopicNotFound.Error())
		assert.Equal(t, "topic already exists", ErrTopicAlreadyExists.Error())
		assert.Equal(t, "topic queue is full", ErrTopicQueueFull.Error())
	})
}
