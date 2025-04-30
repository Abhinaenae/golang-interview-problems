// This queue implementation uses a mutex-based approach using a slice instead of
// managing increased complexity to accomodate the new peek() function that reads values
// without removing them from the store

package main

import (
	"errors"
	"sync"
)

var ErrQueueFull = errors.New("queue is full")

type Queue struct {
	store   []int
	maxSize int
	mu      sync.Mutex
}

func NewQueue(size int) *Queue {
	return &Queue{
		store:   make([]int, 0, size),
		maxSize: size,
		mu:      sync.Mutex{},
	}
}

func (q *Queue) Push(val int) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.store) >= q.maxSize {
		return ErrQueueFull
	}

	q.store = append(q.store, val)
	return nil
}

func (q *Queue) Pop() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.store) == 0 {
		return -1
	}

	poppedVal := q.store[0]
	q.store = q.store[1:]
	return poppedVal
}

func (q *Queue) Peek() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.store) == 0 {
		return -1
	}

	peekVal := q.store[0]
	return peekVal
}
