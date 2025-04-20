package main

import (
	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	defer close(ch)
	walkTree(t, ch)
}

func walkTree(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}
	walkTree(t.Left, ch)
	ch <- t.Value
	walkTree(t.Right, ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	// If both trees are nil, they are the same
	if t1 == nil && t2 == nil {
		return true
	}

	// If one of the trees is nil and the other is not, they are not the same
	if t1 == nil || t2 == nil {
		return false
	}

	// If both trees are not nil, we need to compare their values
	ch1, ch2 := make(chan int, 10), make(chan int, 10) // Buffered channels to avoid deadlocks
	go Walk(t1, ch1)
	go Walk(t2, ch2)
	for v1 := range ch1 {
		v2, ok := <-ch2
		if !ok || v1 != v2 {
			return false
		}
	}

	// If ch2 is closed, both trees are the same
	// If ch2 is not closed, it means t2 has more values than t1
	// and they are not the same
	_, ok := <-ch2
	return !ok
}
