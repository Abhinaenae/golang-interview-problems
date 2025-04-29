package main

import (
	"errors"
)

// For this specific problem, the channel-based implementation is better because
// The code is more concise and has fewer moving parts
// Channels are designed specifically for concurrent communication
// Channel operations are optimized for concurrent access
// The problem requires only Push and Pop operations, which channels handle well
// Channel's buffer naturally implements the fixed-size requirement
// --------------------------------------------------------------
// The mutex-based implementation would be better if:
// You needed to iterate over queue contents
// You needed to peek at elements
// You needed more complex queue operations
// You needed to dynamically resize the queue
// Since none of these are required by the problem specification,
// the channel-based implementation is the superior choice.

var ErrQueueFull = errors.New("queue is full")

type Queue struct {
	store chan int
}

func NewQueue(size int) *Queue {
	return &Queue{store: make(chan int, size)}
}

func (q *Queue) Push(val int) error {
	select {
	case q.store <- val:
		return nil
	default:
		return ErrQueueFull
	}
}

func (q *Queue) Pop() int {
	select {
	case val := <-q.store:
		return val
	default:
		return -1
	}
}
