// A zero allocation queue is a queue that doesn't make any heap allocations during its operations
// Uses a circular buffer approach
package main

import (
	"errors"
	"sync"
)

var ErrQueueFull = errors.New("queue is full")

type Queue struct {
	data []int
	head int //index of first element
	tail int //index where next element will be inserted
	size int //current number of elements
	cap  int // maximum capacity of elements
	mu   sync.Mutex
}

func NewQueue(size int) *Queue {
	return &Queue{
		data: make([]int, size),
		cap:  size,
	}
}

func (q *Queue) Push(val int) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.size >= q.cap {
		return ErrQueueFull
	}

	q.data[q.tail] = val
	q.tail = (q.tail + 1) % q.cap
	q.size++
	return nil
}

func (q *Queue) Pop() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.size == 0 {
		return -1
	}

	val := q.data[q.head]
	q.head = (q.head + 1) % q.cap
	q.size--
	return val
}

func (q *Queue) Peek() int {

	q.mu.Lock()
	defer q.mu.Unlock()

	if q.size == 0 {
		return -1
	}
	return q.data[q.head]
}
