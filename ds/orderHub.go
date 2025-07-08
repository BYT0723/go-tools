package ds

import (
	"sync"
)

var _ Hub[int] = (*OrderHub[int])(nil)

// OrderHub is a simple in-memory publish-subscribe hub
// that delivers messages to subscribers in the order they subscribed.
// Unlike FastBus, OrderHub maintains a slice to preserve subscription order.
type OrderHub[T any] struct {
	// bufSize sets the buffer capacity for each subscriber channel.
	bufSize int

	// subscribers holds all active subscriber channels in subscription order.
	subscribers []chan T

	mutex sync.Mutex
	closeModule
}

// NewOrderHub creates a new OrderBus with the given buffer size for each subscriber.
func NewOrderHub[T any](bufSize int) *OrderHub[T] {
	return &OrderHub[T]{
		bufSize:     bufSize,
		subscribers: make([]chan T, 0, 256),
	}
}

// Subscribe registers a new subscriber and returns its Subscription.
// Messages will be delivered in the order subscribers were added.
// If the bus is closed, it returns ErrHubClosed.
func (b *OrderHub[T]) Subscribe() (*Subscription[T], error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if err := b.closeCheck(); err != nil {
		return nil, err
	}

	ch := make(chan T, b.bufSize)
	s := &Subscription[T]{C: ch}
	b.subscribers = append(b.subscribers, ch)
	return s, nil
}

// Publish sends the given value to all subscribers in subscription order.
// If a subscriber's channel is full, the message is dropped for that subscriber.
// If the bus is closed, it returns ErrHubClosed.
func (b *OrderHub[T]) Publish(v T) error {
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
// After removal, the subscriber will no longer receive messages.
func (b *OrderHub[T]) Unsubscribe(s *Subscription[T]) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for i, ch := range b.subscribers {
		if ch == s.C {
			close(ch)
			b.subscribers = append(b.subscribers[:i], b.subscribers[i+1:]...)
			break
		}
	}
}

// Close closes the bus and all subscriber channels.
// After closing, Subscribe and Publish will return ErrHubClosed.
func (b *OrderHub[T]) Close() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if err := b.closeCheck(); err != nil {
		return err
	}

	for _, ch := range b.subscribers {
		close(ch)
	}
	b.subscribers = nil
	b.closed = true
	return nil
}
