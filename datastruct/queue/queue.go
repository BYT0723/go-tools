package queue

// Queue 定义一个泛型队列接口
type Queue[T any] interface {
	// 入队
	Enqueue(value T)

	// 出队
	Dequeue() (T, bool)

	// 查询
	Front() (T, bool)
	IsEmpty() bool
	Len() int

	// 其他
	Clear()
}
