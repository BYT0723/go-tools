package ds

import (
	"sync"
)

// BroadcastBus is a multi-topic, multi-subscriber pub-sub bus.
//
// Each topic is a FastHub (fan-out):
// - When you publish a message to a topic, all subscribers receive a copy.
// - If a topic has no subscribers, messages are discarded.
//
// Typical use case: event broadcast, notifications, state updates.
type BroadcastBus[T any] struct {
	mutex       sync.Mutex             // protects subscribers map
	subscribers map[string]*FastHub[T] // topic -> FastHub
	closeModule                        // handles closed state
}

// NewBroadcastBus creates a new BroadcastBus with the given buffer size
// for each subscriber channel.
func NewBroadcastBus[T any]() *BroadcastBus[T] {
	return &BroadcastBus[T]{
		subscribers: make(map[string]*FastHub[T]),
	}
}

// AddTopic creates a new topic with its own buffered channel.
// Returns an error if the topic already exists.
func (u *BroadcastBus[T]) AddTopic(topic string, bufSize int) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if err := u.closeCheck(); err != nil {
		return err
	}

	if _, ok := u.subscribers[topic]; ok {
		return ErrTopicAlreadyExists
	}
	u.subscribers[topic] = NewFastHub[T](bufSize)
	return nil
}

// RemoveTopic closes and removes the topic's channel.
// Messages in the buffer are dropped.
func (u *BroadcastBus[T]) RemoveTopic(topic string) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if err := u.closeCheck(); err != nil {
		return err
	}

	if ch, ok := u.subscribers[topic]; ok {
		delete(u.subscribers, topic)
		return ch.Close()
	}
	return ErrTopicNotFound
}

// Subscribe subscribes to the given topic, creating the topic if needed.
//
// Returns a new Subscription which has a buffered channel for receiving messages.
//
// If the bus is closed, returns an error.
func (u *BroadcastBus[T]) Subscribe(topic string) (*Subscription[T], error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if err := u.closeCheck(); err != nil {
		return nil, err
	}

	hub, ok := u.subscribers[topic]
	if !ok {
		return nil, ErrTopicNotFound
	}
	return hub.Subscribe()
}

// Unsubscribe removes a subscriber from the given topic.
// If no subscribers remain for that topic, the topic is closed and removed.
//
// Returns an error if the bus is closed or the topic does not exist.
func (u *BroadcastBus[T]) Unsubscribe(topic string, sub *Subscription[T]) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if err := u.closeCheck(); err != nil {
		return err
	}

	hub, ok := u.subscribers[topic]
	if !ok {
		return ErrTopicNotFound
	}

	hub.Unsubscribe(sub)
	return nil
}

// Publish broadcasts a message to all subscribers of the given topic.
//
// If the topic has no subscribers, returns ErrTopicNotFound.
// If the bus is closed, returns an error.
func (u *BroadcastBus[T]) Publish(topic string, msg T) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if err := u.closeCheck(); err != nil {
		return err
	}

	hub, ok := u.subscribers[topic]
	if !ok {
		return ErrTopicNotFound
	}

	return hub.Publish(msg)
}

// Close closes all topics and marks the bus as closed.
// Any further operations on the bus will return an error.
func (u *BroadcastBus[T]) Close() error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if err := u.closeCheck(); err != nil {
		return err
	}

	for _, hub := range u.subscribers {
		if err := hub.Close(); err != nil {
			return err
		}
	}

	u.subscribers = nil
	u.closed = true
	return nil
}
