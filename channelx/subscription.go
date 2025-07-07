package channelx

// Subscription represents a single subscriber to a Hub.
// It holds the unique subscription ID and the receive-only channel
// that delivers published messages.
type Subscription[T any] struct {
	// id is the unique identifier for this subscriber.
	id uint64

	// C is the receive-only channel for published messages.
	C <-chan T
}
