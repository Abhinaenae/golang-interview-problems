package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Client interface {
	Get(ctx context.Context, address string) (string, error)
}

var ErrRequestsFailed = errors.New("task execution failed")

// RequestWithFailover attempts to request a data from available addresses:
// 1. If error, immediately try the address without waiting
// 2. If an address doesn't respond within 500ms, try the next but keep the original request running
// 3. Return the first successful response, or all ErrRequestsFailed if all nodes fail
// 4. Properly handle context cancellation throughout the process
func RequestWithFailover(ctx context.Context, client Client, addresses []string) (string, error) {
	resCh := make(chan string, 1)
	errCh := make(chan error, len(addresses))
	var wg sync.WaitGroup
	var once sync.Once

	for _, address := range addresses {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()

			// Create timeout context only for failover mechanism
			timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			defer timeoutCancel()

			// Launch the actual request with parent context
			done := make(chan struct{})
			go func() {
				res, err := client.Get(ctx, addr)
				if err != nil {
					errCh <- err
					return
				}

				once.Do(func() {
					select {
					case resCh <- res:
						fmt.Printf("Successful response from %s: %s\n", addr, res)
					default:
					}
				})
				close(done)
			}()

			// Wait for either timeout or completion
			select {
			case <-timeoutCtx.Done():
				// Allow request to continue in background
				return
			case <-done:
				// Request completed before timeout
				return
			}
		}(address)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	var errCnt int
	for {
		select {
		case resp := <-resCh:
			return resp, nil
		case _, ok := <-errCh:
			if !ok {
				// errCh closed and no successful response
				if errCnt == len(addresses) {
					return "", ErrRequestsFailed
				}
				continue
			}
			errCnt++
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}
