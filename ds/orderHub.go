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
	subscribers []*Subscription[T]

	mutex sync.Mutex
	closeModule
	idGenerator
	hubCallback[T]
}

// NewOrderHub creates a new OrderBus with the given buffer size for each subscriber.
func NewOrderHub[T any](bufSize int) *OrderHub[T] {
	return &OrderHub[T]{
		bufSize:     bufSize,
		subscribers: make([]*Subscription[T], 0, 256),
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
	s := &Subscription[T]{
		ID: b.Increment(),
		c:  ch,
	}
	b.subscribers = append(b.subscribers, s)
	return s, nil
}

// Publish sends the given value to all subscribers in subscription order.
// If a subscriber's channel is full, the message is dropped for that subscriber.
// If the bus is closed, it returns ErrHubClosed.
func (b *OrderHub[T]) Publish(v T) error {
	b.mutex.Lock()
	if err := b.closeCheck(); err != nil {
		b.mutex.Unlock()
		return err
	}
	subs := append([]*Subscription[T](nil), b.subscribers...)
	b.mutex.Unlock()

	for _, sub := range subs {
		select {
		case sub.c <- v:
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

	for i, sub := range b.subscribers {
		if sub == s {
			close(sub.c)
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

	for _, sub := range b.subscribers {
		close(sub.c)
	}
	b.subscribers = nil
	b.closed = true
	return nil
}
