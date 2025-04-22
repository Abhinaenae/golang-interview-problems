package main

import "sync"

type Client interface {
	Get(address string) (string, error)
}

// Cache is a non-blocking cache that caches the result of a Get call
// It uses a map to store the results and a mutex to protect access to the map
// It uses a channel to signal when the result is ready
type Cache struct {
	client  Client
	m       map[string]*data
	mapLock sync.Mutex
}

// data is a struct that holds the result of a Get call
// It uses a channel to signal when the result is ready
// It uses a string to hold the result of the Get call
// It uses an error to hold the error of the Get call
type data struct {
	body  string
	err   error
	ready chan struct{}
}

// NewCache creates a new Cache
// It takes a Client as an argument
// It returns a pointer to a Cache
func NewCache(client Client) *Cache {
	return &Cache{
		client: client,
		m:      make(map[string]*data, 10),
	}
}

// Cache Client.Get result
// This pattern is commonly used to prevent "thundering herd" problems in distributed systems,
// where multiple concurrent requests for the same resource could overwhelm the system.
func (c *Cache) Get(address string) (string, error) {
	c.mapLock.Lock()
	dataRetrieved, ok := c.m[address]
	if !ok {
		dataRetrieved = &data{
			body:  "",
			err:   nil,
			ready: make(chan struct{}),
		}

		c.m[address] = dataRetrieved
		c.mapLock.Unlock()

		dataRetrieved.body, dataRetrieved.err = c.client.Get(address)
		close(dataRetrieved.ready)
	} else {
		c.mapLock.Unlock()
		<-dataRetrieved.ready
	}
	return dataRetrieved.body, dataRetrieved.err
}
