package ds

import (
	"errors"
)

// ErrHubClosed is returned when operations are attempted on a closed Hub.
var (
	ErrHubClosed = errors.New("hub is closed")
)

// Hub defines a simple publish-subscribe broadcast center.
// All subscribers receive every published message.
// It does not guarantee message ordering or delivery reliability.
type Hub[T any] interface {
	// Subscribe registers a new subscriber and returns its Subscription.
	// If the Hub is closed, it returns ErrHubClosed.
	Subscribe() (*Subscription[T], error)

	// Publish sends the given value to all active subscribers.
	// If the Hub is closed, it returns ErrHubClosed.
	Publish(v T) error

	// Unsubscribe removes the given Subscription from the Hub.
	// After unsubscription, the subscriber will no longer receive messages.
	Unsubscribe(*Subscription[T])

	// SetPublishCallback sets a custom publish callback function.
	// callback is called asynchronously when a message cannot be delivered to a subscriber due to a full buffer.
	// The callback must be non-blocking and must NOT call back into Hub.
	SetPublishCallback(callback func(id uint64, v T))

	// Close closes the Hub and all active subscriptions.
	// After closing, Subscribe and Publish will return ErrHubClosed.
	Close() error
}

// Subscription represents a single subscriber to a Hub.
// It holds the unique subscription ID and the receive-only channel
// that delivers published messages.
type Subscription[T any] struct {
	// ID is the unique identifier for this subscriber.
	ID uint64
	// C is the receive-only channel for published messages.
	c chan T
}

// Channel returns the receive-only channel for this subscription.
// Use this channel to receive messages published to the Hub.
func (s *Subscription[T]) Channel() <-chan T {
	return s.c
}

// closeModule is an embedded helper for implementing
// the closed-state check in a Hub implementation.
type closeModule struct {
	closed bool
}

// closeCheck checks if the Hub is closed.
// It returns ErrHubClosed if closed, otherwise nil.
func (m *closeModule) closeCheck() error {
	if m.closed {
		return ErrHubClosed
	}
	return nil
}

type hubCallback[T any] struct {
	publishCallback func(id uint64, v T)
}

func (h *hubCallback[T]) SetPublishCallback(f func(id uint64, v T)) {
	h.publishCallback = f
}
