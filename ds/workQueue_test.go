package ds

import (
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWorkQueueAddTopic(t *testing.T) {
	Convey("WorkQueue AddTopic 测试", t, func() {
		wq := NewWorkQueue[int]()

		Convey("添加新topic成功", func() {
			err := wq.AddTopic("topic1", 10)
			So(err, ShouldBeNil)
		})

		Convey("重复添加同一topic返回错误", func() {
			wq.AddTopic("topic1", 10)
			err := wq.AddTopic("topic1", 10)
			So(err, ShouldEqual, ErrTopicAlreadyExists)
		})
	})
}

func TestWorkQueueRemoveTopic(t *testing.T) {
	Convey("WorkQueue RemoveTopic 测试", t, func() {
		Convey("移除已存在的topic", func() {
			wq := NewWorkQueue[int]()
			wq.AddTopic("topic1", 10)
			err := wq.RemoveTopic("topic1")
			So(err, ShouldBeNil)
		})

		Convey("移除不存在的topic不返回错误", func() {
			wq := NewWorkQueue[int]()
			err := wq.RemoveTopic("nonexistent")
			So(err, ShouldBeNil)
		})
	})
}

func TestWorkQueuePublish(t *testing.T) {
	Convey("WorkQueue Publish 测试", t, func() {
		Convey("发布消息到已有topic", func() {
			wq := NewWorkQueue[int]()
			wq.AddTopic("topic1", 10)
			err := wq.Publish("topic1", 42)
			So(err, ShouldBeNil)
		})

		Convey("发布消息到不存在topic返回错误", func() {
			wq := NewWorkQueue[int]()
			err := wq.Publish("nonexistent", 42)
			So(err, ShouldEqual, ErrTopicNotFound)
		})

		Convey("发布到满队列返回队列满错误", func() {
			wq := NewWorkQueue[int]()
			wq.AddTopic("topic1", 1)
			wq.Publish("topic1", 1)
			err := wq.Publish("topic1", 2)
			So(err, ShouldEqual, ErrTopicQueueFull)
		})
	})
}

func TestWorkQueueSubscribe(t *testing.T) {
	Convey("WorkQueue Subscribe 测试", t, func() {
		Convey("订阅已存在的topic成功", func() {
			wq := NewWorkQueue[int]()
			wq.AddTopic("topic1", 10)
			ch, err := wq.Subscribe("topic1")
			So(err, ShouldBeNil)
			So(ch, ShouldNotBeNil)
		})

		Convey("订阅不存在的topic返回错误", func() {
			wq := NewWorkQueue[int]()
			ch, err := wq.Subscribe("nonexistent")
			So(err, ShouldEqual, ErrTopicNotFound)
			So(ch, ShouldBeNil)
		})
	})
}

func TestWorkQueuePubSub(t *testing.T) {
	Convey("WorkQueue 发布订阅消息传递测试", t, func() {
		Convey("订阅者能收到发布的消息", func() {
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
			So(result, ShouldEqual, 42)
		})

		Convey("两个topic互不影响", func() {
			wq := NewWorkQueue[int]()
			wq.AddTopic("topic1", 10)
			wq.AddTopic("topic2", 10)

			ch1, _ := wq.Subscribe("topic1")
			ch2, _ := wq.Subscribe("topic2")

			wq.Publish("topic1", 1)
			wq.Publish("topic2", 2)

			So(<-ch1, ShouldEqual, 1)
			So(<-ch2, ShouldEqual, 2)
		})
	})
}

func TestWorkQueueClose(t *testing.T) {
	Convey("WorkQueue Close 测试", t, func() {
		Convey("Close 后channel被关闭", func() {
			wq := NewWorkQueue[int]()
			wq.AddTopic("topic1", 10)
			wq.Close()

			ch, _ := wq.Subscribe("topic1")
			_, ok := <-ch
			So(ok, ShouldBeFalse)
		})
	})
}

func TestWorkQueueErrors(t *testing.T) {
	Convey("WorkQueue 错误变量测试", t, func() {
		So(ErrTopicNotFound.Error(), ShouldEqual, "topic not found")
		So(ErrTopicAlreadyExists.Error(), ShouldEqual, "topic already exists")
		So(ErrTopicQueueFull.Error(), ShouldEqual, "topic queue is full")
	})
}
