package ds

// Map is a generic interface for thread-safe key-value storage.
// It provides operations similar to Go's sync.Map but with a more comprehensive API.
//
// Type parameters:
//   - K: The key type, must be comparable
//   - V: The value type
//
// Implementations:
//   - SyncMap: Wraps Go's sync.Map with type safety
//   - MutexMap: Uses a mutex-protected map for simpler concurrency control
type Map[K comparable, V any] interface {
	// Store stores a value for a key.
	Store(K, V)

	// Load returns the value stored for a key, or false if no value is present.
	Load(K) (V, bool)

	// Delete deletes the value for a key.
	Delete(K) bool

	// Swap swaps the value for a key and returns the previous value if any.
	// The loaded result reports whether the key was present.
	Swap(key K, newValue V) (old V, loaded bool)

	// Range calls iterator sequentially for each key and value present in the map.
	// If iterator returns false, range stops the iteration.
	Range(iterator func(K, V) bool)

	// LoadOrStore returns the existing value for the key if present.
	// Otherwise, it stores and returns the given value.
	// The loaded result is true if the value was loaded, false if stored.
	LoadOrStore(key K, newValue V) (value V, loaded bool)

	// LoadAndDelete deletes the value for a key, returning the previous value if any.
	// The loaded result reports whether the key was present.
	LoadAndDelete(K) (value V, loaded bool)

	// CompareAndSwap swaps the value for a key if the current value equals old.
	// Returns true if the swap was performed.
	CompareAndSwap(key K, old, newValue V) bool

	// CompareAndDelete deletes the entry for a key if its value equals old.
	// Returns true if the entry was deleted.
	CompareAndDelete(key K, value V) bool

	// CompareFnAndSwap swaps the value for a key using a custom comparison function.
	// The function fn is called with the current value and old value.
	// Returns true if fn returns true and the swap was performed.
	CompareFnAndSwap(key K, fn func(V, V) bool, old, newValue V) bool

	// CompareFnAndDelete deletes the entry for a key using a custom comparison function.
	// The function fn is called with the current value and old value.
	// Returns true if fn returns true and the entry was deleted.
	CompareFnAndDelete(key K, fn func(V, V) bool, old V) bool

	// Keys returns a slice containing all keys in the map.
	Keys() []K

	// Values returns a slice containing all values in the map.
	Values() []V

	// Filter returns a new Map containing only entries that satisfy the filter function.
	Filter(filter func(K, V) bool) Map[K, V]
}
