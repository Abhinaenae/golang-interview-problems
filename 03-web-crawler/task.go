package main

import "sync"

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type store struct {
	visited map[string]bool
	mu      sync.RWMutex
}

var storeMap = store{
	visited: make(map[string]bool),
	mu:      sync.RWMutex{},
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) ([]string, error) {
	storeMap.mu.RLock()
	if storeMap.visited[url] {
		storeMap.mu.RUnlock()
		return nil, nil
	}
	storeMap.mu.RUnlock()

	storeMap.mu.Lock()
	storeMap.visited[url] = true
	storeMap.mu.Unlock()
	if depth <= 0 {
		return nil, nil
	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		return nil, err
	}

	result := []string{body}

	if len(urls) == 0 {
		return result, nil
	}

	//Use waitgroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	for _, u := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if res, err := Crawl(u, depth-1, fetcher); err == nil {
				result = append(result, res...)
			}
		}()
	}
	// Wait for all goroutines to finish
	wg.Wait()
	return result, nil
}
