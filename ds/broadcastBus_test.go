package ds

import (
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBroadcastBusAddTopic(t *testing.T) {
	Convey("BroadcastBus AddTopic 测试", t, func() {
		bus := NewBroadcastBus[int]()

		Convey("添加新topic成功", func() {
			err := bus.AddTopic("topic1", 10)
			So(err, ShouldBeNil)
		})

		Convey("重复添加同一topic返回错误", func() {
			bus.AddTopic("topic1", 10)
			err := bus.AddTopic("topic1", 10)
			So(err, ShouldEqual, ErrTopicAlreadyExists)
		})

		Convey("关闭后添加topic返回错误", func() {
			bus.Close()
			err := bus.AddTopic("topic1", 10)
			So(err, ShouldEqual, ErrHubClosed)
		})
	})
}

func TestBroadcastBusRemoveTopic(t *testing.T) {
	Convey("BroadcastBus RemoveTopic 测试", t, func() {
		Convey("移除已存在的topic成功", func() {
			bus := NewBroadcastBus[int]()
			bus.AddTopic("topic1", 10)
			err := bus.RemoveTopic("topic1")
			So(err, ShouldBeNil)
		})

		Convey("移除不存在的topic返回错误", func() {
			bus := NewBroadcastBus[int]()
			err := bus.RemoveTopic("nonexistent")
			So(err, ShouldEqual, ErrTopicNotFound)
		})
	})
}

func TestBroadcastBusSubscribe(t *testing.T) {
	Convey("BroadcastBus Subscribe 测试", t, func() {
		Convey("订阅已存在的topic成功", func() {
			bus := NewBroadcastBus[int]()
			bus.AddTopic("topic1", 10)
			sub, err := bus.Subscribe("topic1")
			So(err, ShouldBeNil)
			So(sub, ShouldNotBeNil)
		})

		Convey("订阅不存在的topic返回错误", func() {
			bus := NewBroadcastBus[int]()
			sub, err := bus.Subscribe("nonexistent")
			So(err, ShouldEqual, ErrTopicNotFound)
			So(sub, ShouldBeNil)
		})

		Convey("关闭后订阅返回错误", func() {
			bus := NewBroadcastBus[int]()
			bus.AddTopic("topic1", 10)
			bus.Close()
			sub, err := bus.Subscribe("topic1")
			So(err, ShouldEqual, ErrHubClosed)
			So(sub, ShouldBeNil)
		})
	})
}

func TestBroadcastBusUnsubscribe(t *testing.T) {
	Convey("BroadcastBus Unsubscribe 测试", t, func() {
		Convey("取消订阅成功", func() {
			bus := NewBroadcastBus[int]()
			bus.AddTopic("topic1", 10)
			sub, _ := bus.Subscribe("topic1")
			err := bus.Unsubscribe("topic1", sub)
			So(err, ShouldBeNil)
		})

		Convey("取消订阅不存在的topic返回错误", func() {
			bus := NewBroadcastBus[int]()
			bus.AddTopic("topic1", 10)
			sub, _ := bus.Subscribe("topic1")
			err := bus.Unsubscribe("nonexistent", sub)
			So(err, ShouldEqual, ErrTopicNotFound)
		})
	})
}

func TestBroadcastBusPublish(t *testing.T) {
	Convey("BroadcastBus Publish 测试", t, func() {
		Convey("发布消息到已有topic", func() {
			bus := NewBroadcastBus[int]()
			bus.AddTopic("topic1", 10)
			err := bus.Publish("topic1", 42)
			So(err, ShouldBeNil)
		})

		Convey("发布消息到不存在topic返回错误", func() {
			bus := NewBroadcastBus[int]()
			err := bus.Publish("nonexistent", 42)
			So(err, ShouldEqual, ErrTopicNotFound)
		})

		Convey("关闭后发布返回错误", func() {
			bus := NewBroadcastBus[int]()
			bus.AddTopic("topic1", 10)
			bus.Close()
			err := bus.Publish("topic1", 42)
			So(err, ShouldEqual, ErrHubClosed)
		})
	})
}

func TestBroadcastBusPublishReceive(t *testing.T) {
	Convey("BroadcastBus 发布订阅消息传递测试", t, func() {
		Convey("多订阅者收到相同消息", func() {
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
				So(msg, ShouldEqual, "hello")
			}
		})

		Convey("单订阅者单消息", func() {
			bus := NewBroadcastBus[int]()
			bus.AddTopic("topic1", 1)
			sub, err := bus.Subscribe("topic1")
			So(err, ShouldBeNil)

			bus.Publish("topic1", 42)
			msg := <-sub.Channel()
			So(msg, ShouldEqual, 42)
		})
	})
}

func TestBroadcastBusClose(t *testing.T) {
	Convey("BroadcastBus Close 测试", t, func() {
		Convey("正常关闭", func() {
			bus := NewBroadcastBus[int]()
			bus.AddTopic("topic1", 10)
			bus.AddTopic("topic2", 10)
			err := bus.Close()
			So(err, ShouldBeNil)
		})

		Convey("重复关闭返回错误", func() {
			bus := NewBroadcastBus[int]()
			bus.Close()
			err := bus.Close()
			So(err, ShouldEqual, ErrHubClosed)
		})
	})
}
