package main

import (
	"context"
	"errors"
	"strings"
	"sync"
)

type Getter interface {
	Get(ctx context.Context, address, key string) (string, error)
}

// Call `Getter.Get()` for each address in parallel.
// Returns the first successful response.
// If all requests fail, returns an error.
func Get(ctx context.Context, getter Getter, addresses []string, key string) (string, error) {

	if len(addresses) == 0 {
		return "", nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Channels MUST be buffered, in other case there is a goroutine leakage
	resCh := make(chan result, len(addresses))
	var wg sync.WaitGroup

	for _, address := range addresses {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			val, err := getter.Get(ctx, addr, key)
			select {
			case resCh <- result{value: val, err: err}:
			case <-ctx.Done():
			}
		}(address)
	}

	// Wait for all goroutines to finish and close the result channel
	go func() {
		wg.Wait()
		close(resCh)
	}()

	var errs []string
	for res := range resCh {
		if res.err == nil {
			cancel() //Cancel other goroutines once have a successful result
			return res.value, nil
		}
		errs = append(errs, res.err.Error())

	}

	//Aggregate all errors into single error message
	return "", errors.New("all requests failed: " + strings.Join(errs, "; "))
}

type result struct {
	value string
	err   error
}
