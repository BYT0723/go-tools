package ds

import (
	"errors"
	"sync"
)

// WorkQueue is a simple pub-sub style queue where each topic
// has exactly one subscriber (worker).
// A published message goes to exactly one consumer.
//
// Typical use case: task queues, job workers.
//
// T is the type of the message.
type WorkQueue[T any] struct {
	mutex       sync.Mutex        // protects subscribers map
	subscribers map[string]chan T // topic name -> channel
}

var (
	// ErrTopicNotFound indicates the given topic does not exist.
	ErrTopicNotFound = errors.New("topic not found")

	// ErrTopicAlreadyExists indicates the topic already exists.
	ErrTopicAlreadyExists = errors.New("topic already exists")

	// ErrTopicQueueFull indicates the topic queue is full
	ErrTopicQueueFull = errors.New("topic queue is full")
)

// NewWorkQueue creates a new WorkQueue with a given channel buffer size.
func NewWorkQueue[T any]() *WorkQueue[T] {
	return &WorkQueue[T]{
		subscribers: make(map[string]chan T),
	}
}

// AddTopic creates a new topic with its own buffered channel.
// Returns an error if the topic already exists.
func (u *WorkQueue[T]) AddTopic(topic string, bufSize int) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if _, ok := u.subscribers[topic]; ok {
		return ErrTopicAlreadyExists
	}
	u.subscribers[topic] = make(chan T, bufSize)
	return nil
}

// RemoveTopic closes and removes the topic's channel.
// Messages in the buffer are dropped.
func (u *WorkQueue[T]) RemoveTopic(topic string) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if ch, ok := u.subscribers[topic]; ok {
		close(ch)
		delete(u.subscribers, topic)
	}
	return nil
}

// Publish sends a message to the topic's channel.
// If the channel buffer is full, the message is dropped silently.
// Returns an error if the topic does not exist.
func (u *WorkQueue[T]) Publish(topic string, msg T) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if ch, ok := u.subscribers[topic]; ok {
		select {
		case ch <- msg:
		default:
			return ErrTopicQueueFull
		}
		return nil
	}
	return ErrTopicNotFound
}

// Subscribe returns the channel for the given topic.
// Only one subscriber should consume from the channel.
// Returns an error if the topic does not exist.
func (u *WorkQueue[T]) Subscribe(topic string) (<-chan T, error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if ch, ok := u.subscribers[topic]; ok {
		return ch, nil
	}
	return nil, ErrTopicNotFound
}

// Close closes all topic channels.
// Topics are not removed from the map.
func (u *WorkQueue[T]) Close() {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	for _, ch := range u.subscribers {
		close(ch)
	}
}
