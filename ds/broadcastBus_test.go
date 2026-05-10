package ds

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBroadcastBusAddTopic(t *testing.T) {
	t.Run("BroadcastBus AddTopic 测试", func(t *testing.T) {
		bus := NewBroadcastBus[int]()

		t.Run("添加新topic成功", func(t *testing.T) {
			err := bus.AddTopic("topic1", 10)
			assert.Nil(t, err)
		})

		t.Run("重复添加同一topic返回错误", func(t *testing.T) {
			bus.AddTopic("topic1", 10)
			err := bus.AddTopic("topic1", 10)
			assert.Equal(t, ErrTopicAlreadyExists, err)
		})

		t.Run("关闭后添加topic返回错误", func(t *testing.T) {
			bus.Close()
			err := bus.AddTopic("topic1", 10)
			assert.Equal(t, ErrHubClosed, err)
		})
	})
}

func TestBroadcastBusRemoveTopic(t *testing.T) {
	t.Run("BroadcastBus RemoveTopic 测试", func(t *testing.T) {
		t.Run("移除已存在的topic成功", func(t *testing.T) {
			bus := NewBroadcastBus[int]()
			bus.AddTopic("topic1", 10)
			err := bus.RemoveTopic("topic1")
			assert.Nil(t, err)
		})

		t.Run("移除不存在的topic返回错误", func(t *testing.T) {
			bus := NewBroadcastBus[int]()
			err := bus.RemoveTopic("nonexistent")
			assert.Equal(t, ErrTopicNotFound, err)
		})
	})
}

func TestBroadcastBusSubscribe(t *testing.T) {
	t.Run("BroadcastBus Subscribe 测试", func(t *testing.T) {
		t.Run("订阅已存在的topic成功", func(t *testing.T) {
			bus := NewBroadcastBus[int]()
			bus.AddTopic("topic1", 10)
			sub, err := bus.Subscribe("topic1")
			assert.Nil(t, err)
			assert.NotNil(t, sub)
		})

		t.Run("订阅不存在的topic返回错误", func(t *testing.T) {
			bus := NewBroadcastBus[int]()
			sub, err := bus.Subscribe("nonexistent")
			assert.Equal(t, ErrTopicNotFound, err)
			assert.Nil(t, sub)
		})

		t.Run("关闭后订阅返回错误", func(t *testing.T) {
			bus := NewBroadcastBus[int]()
			bus.AddTopic("topic1", 10)
			bus.Close()
			sub, err := bus.Subscribe("topic1")
			assert.Equal(t, ErrHubClosed, err)
			assert.Nil(t, sub)
		})
	})
}

func TestBroadcastBusUnsubscribe(t *testing.T) {
	t.Run("BroadcastBus Unsubscribe 测试", func(t *testing.T) {
		t.Run("取消订阅成功", func(t *testing.T) {
			bus := NewBroadcastBus[int]()
			bus.AddTopic("topic1", 10)
			sub, _ := bus.Subscribe("topic1")
			err := bus.Unsubscribe("topic1", sub)
			assert.Nil(t, err)
		})

		t.Run("取消订阅不存在的topic返回错误", func(t *testing.T) {
			bus := NewBroadcastBus[int]()
			bus.AddTopic("topic1", 10)
			sub, _ := bus.Subscribe("topic1")
			err := bus.Unsubscribe("nonexistent", sub)
			assert.Equal(t, ErrTopicNotFound, err)
		})
	})
}

func TestBroadcastBusPublish(t *testing.T) {
	t.Run("BroadcastBus Publish 测试", func(t *testing.T) {
		t.Run("发布消息到已有topic", func(t *testing.T) {
			bus := NewBroadcastBus[int]()
			bus.AddTopic("topic1", 10)
			err := bus.Publish("topic1", 42)
			assert.Nil(t, err)
		})

		t.Run("发布消息到不存在topic返回错误", func(t *testing.T) {
			bus := NewBroadcastBus[int]()
			err := bus.Publish("nonexistent", 42)
			assert.Equal(t, ErrTopicNotFound, err)
		})

		t.Run("关闭后发布返回错误", func(t *testing.T) {
			bus := NewBroadcastBus[int]()
			bus.AddTopic("topic1", 10)
			bus.Close()
			err := bus.Publish("topic1", 42)
			assert.Equal(t, ErrHubClosed, err)
		})
	})
}

func TestBroadcastBusPublishReceive(t *testing.T) {
	t.Run("BroadcastBus 发布订阅消息传递测试", func(t *testing.T) {
		t.Run("多订阅者收到相同消息", func(t *testing.T) {
			bus := NewBroadcastBus[string]()
			bus.AddTopic("topic1", 10)

			var wg sync.WaitGroup
			n := 5
			received := make([]string, n)

			for i := 0; i < n; i++ {
				sub, _ := bus.Subscribe("topic1")
				wg.Add(1)
				go func(idx int, s *Subscription[string]) {
					defer wg.Done()
					received[idx] = <-s.Channel()
				}(i, sub)
			}

			bus.Publish("topic1", "hello")

			wg.Wait()

			for _, msg := range received {
				assert.Equal(t, "hello", msg)
			}
		})

		t.Run("单订阅者单消息", func(t *testing.T) {
			bus := NewBroadcastBus[int]()
			bus.AddTopic("topic1", 1)
			sub, err := bus.Subscribe("topic1")
			assert.Nil(t, err)

			bus.Publish("topic1", 42)
			msg := <-sub.Channel()
			assert.Equal(t, 42, msg)
		})
	})
}

func TestBroadcastBusClose(t *testing.T) {
	t.Run("BroadcastBus Close 测试", func(t *testing.T) {
		t.Run("正常关闭", func(t *testing.T) {
			bus := NewBroadcastBus[int]()
			bus.AddTopic("topic1", 10)
			bus.AddTopic("topic2", 10)
			err := bus.Close()
			assert.Nil(t, err)
		})

		t.Run("重复关闭返回错误", func(t *testing.T) {
			bus := NewBroadcastBus[int]()
			bus.Close()
			err := bus.Close()
			assert.Equal(t, ErrHubClosed, err)
		})
	})
}
