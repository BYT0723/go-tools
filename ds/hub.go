package ds

import (
	"errors"
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

	// Close closes the Hub and all active subscriptions.
	// After closing, Subscribe and Publish will return ErrHubClosed.
	Close() error
}

// ErrHubClosed is returned when operations are attempted on a closed Hub.
var ErrHubClosed = errors.New("hub is closed")

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

// Subscription represents a single subscriber to a Hub.
// It holds the unique subscription ID and the receive-only channel
// that delivers published messages.
type Subscription[T any] struct {
	// id is the unique identifier for this subscriber.
	id uint64

	// C is the receive-only channel for published messages.
	C <-chan T
}
