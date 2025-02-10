package ds

type Map[K comparable, V any] interface {
	// 存储
	Store(K, V)
	// 读取
	Load(K) (V, bool)
	// 删除
	Delete(K) bool
	// 交换
	Swap(key K, newValue V) (old V, loaded bool)
	// 遍历，当iterator返回false时，遍历终止
	Range(iterator func(K, V) bool)
	// 读取或将值存储
	LoadOrStore(key K, newValue V) (value V, loaded bool)
	// 读取并删除
	LoadAndDelete(K) (value V, loaded bool)
	// 对比并交换
	// 对比当前map中key对应的value是否等于old，如果相等则将key对应的value替换为newValue
	CompareAndSwap(key K, old, newValue V) bool
	// 对比并删除
	// 对比当前map中key对应的value是否等于value，如果相等则删除
	CompareAndDelete(key K, value V) bool
	CompareFnAndSwap(key K, fn func(V, V) bool, old, newValue V) bool
	CompareFnAndDelete(key K, fn func(V, V) bool, old V) bool
}
