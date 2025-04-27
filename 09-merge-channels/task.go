package main

import (
	"sync"
)

func Merge(channels ...<-chan int) <-chan int {
	out := make(chan int)

	if len(channels) == 0 {
		close(out)
		return out
	}

	// Create WaitGroup for the main processing goroutine
	var wg sync.WaitGroup
	wg.Add(1)

	// Create WaitGroup for the reader goroutines
	var readerWg sync.WaitGroup

	go func() {
		defer wg.Done()

		cases := make([]chan int, len(channels))
		activeChannels := len(channels)

		// Start reader goroutines
		for i := range channels {
			cases[i] = make(chan int)
			readerWg.Add(1)
			go func(inputCh <-chan int, caseCh chan int) {
				defer readerWg.Done()
				defer close(caseCh)
				for val := range inputCh {
					caseCh <- val
				}
			}(channels[i], cases[i])
		}

		// Process values until all channels are closed
		for activeChannels > 0 {
			for i, ch := range cases {
				if ch == nil {
					continue
				}
				select {
				case val, ok := <-ch:
					if !ok {
						cases[i] = nil
						activeChannels--
						continue
					}
					out <- val
				default:
					continue
				}
			}
		}

		// Wait for all readers to finish
		readerWg.Wait()
	}()

	// Close output channel after all processing is done
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
