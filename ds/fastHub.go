package ds

import (
	"sync"
	"sync/atomic"
)

var _ Hub[int] = (*FastHub[int])(nil)

// FastHub is a simple in-memory publish-subscribe hub.
// It broadcasts each published message to all active subscribers.
// FastHub does not guarantee message ordering or reliable delivery.
// If a subscriber's channel buffer is full, the message is dropped.
type FastHub[T any] struct {
	bufSize     int
	subscribers map[uint64]chan T
	mutex       sync.Mutex
	closeModule
	idGenerator
}

// NewFastHub creates a new FastBus with the given buffer size for each subscriber.
func NewFastHub[T any](bufSize int) *FastHub[T] {
	return &FastHub[T]{
		subscribers: make(map[uint64]chan T),
		bufSize:     bufSize,
	}
}

// Subscribe registers a new subscriber and returns its Subscription.
// If the bus is closed, it returns ErrHubClosed.
func (b *FastHub[T]) Subscribe() (*Subscription[T], error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if err := b.closeCheck(); err != nil {
		return nil, err
	}

	var (
		ch = make(chan T, b.bufSize)
		s  = &Subscription[T]{
			id: b.Increment(),
			C:  ch,
		}
	)
	// Ensure ID uniqueness.
	for {
		if _, ok := b.subscribers[s.id]; !ok {
			break
		}
		s.id = b.Increment()
	}
	b.subscribers[s.id] = ch
	return s, nil
}

// Publish broadcasts the given value to all active subscribers.
// If a subscriber's channel is full, the message is dropped.
// If the bus is closed, it returns ErrHubClosed.
func (b *FastHub[T]) Publish(v T) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if err := b.closeCheck(); err != nil {
		return err
	}

	for _, ch := range b.subscribers {
		select {
		case ch <- v:
		default:
			// Drop message if subscriber buffer is full.
		}
	}
	return nil
}

// Unsubscribe removes a subscriber from the bus and closes its channel.
func (b *FastHub[T]) Unsubscribe(s *Subscription[T]) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.closed {
		return
	}

	if ch, ok := b.subscribers[s.id]; ok {
		delete(b.subscribers, s.id)
		close(ch)
	}
}

// Close closes the bus and all active subscriber channels.
// After closing, Publish and Subscribe will return ErrHubClosed.
func (b *FastHub[T]) Close() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.closed {
		return nil
	}

	for _, ch := range b.subscribers {
		close(ch)
	}
	b.subscribers = nil
	b.closed = true
	return nil
}

func (b *FastHub[T]) isEmpty() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return len(b.subscribers) == 0
}

// idGenerator provides atomic ID generation for subscribers.
type idGenerator struct {
	id uint64
}

// Increment atomically increments and returns a new unique ID.
func (g *idGenerator) Increment() uint64 {
	return atomic.AddUint64(&g.id, 1)
}
