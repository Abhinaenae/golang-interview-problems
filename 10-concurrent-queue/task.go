package main

import (
	"errors"
	"sync"
)

var ErrQueueFull = errors.New("queue is full")

type Queue struct {
	store []int
	mu    *sync.Mutex
}

func NewQueue(size int) *Queue {
	return &Queue{
		store: make([]int, 0, size),
		mu:    &sync.Mutex{},
	}
}

func (q *Queue) Push(val int) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.store) == cap(q.store) {
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

	remove := func(slice []int, s int) []int {
		return append(slice[:s], slice[s+1:]...)
	}

	pop := q.store[0]
	q.store = remove(q.store, 0)
	return pop
}
