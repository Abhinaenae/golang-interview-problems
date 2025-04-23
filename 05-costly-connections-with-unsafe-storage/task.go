package main

import (
	"log/slog"
	"sync"
)

type Connection interface {
	// Need call Connect before Send
	// Take time to connect
	Connect()

	// Every connection should be disconnected after use
	// Take time to disconnect
	Disconnect()

	Send(req string) (string, error)
}

type ConnectionCreator interface {
	// Create new connection
	// Will return error if there is more than maxConn
	NewConnection() (Connection, error)
}

type Saver interface {
	// Saves data to unsafe storage
	// WILL CORRUPT DATA on concurrent save
	Save(data string)
}

// SendAndSave should send all requests concurrently using at most `maxConn` simultaneous connections.
// Responses must be saved using Saver.Save.
// Be careful: Saver.Save is not safe for concurrent use.
func SendAndSave(creator ConnectionCreator, saver Saver, requests []string, maxConn int) {
	var wg sync.WaitGroup
	wg.Add(maxConn)
	reqCh := make(chan string, len(requests))
	saveCh := make(chan string, len(requests))

	//Populate request channel
	for _, req := range requests {
		reqCh <- req
	}
	close(reqCh)

	//Connection pool to reuse connections
	connPool := make(chan Connection, maxConn)
	for range maxConn {
		conn, err := creator.NewConnection()
		if err != nil {
			slog.Error("Failed to create connection", "error", err)
			return
		}
		conn.Connect()
		connPool <- conn
	}

	//Worker pool to process requests
	for range maxConn {
		go func() {
			defer wg.Done()

			for req := range reqCh {
				conn := <-connPool //Receive connection from pool
				resp, err := conn.Send(req)
				if err != nil {
					slog.Error("Failed to send request", "error", err)
					connPool <- conn //Send connection back into pool
					continue
				}
				saveCh <- resp
				connPool <- conn //Send connection back into pool
			}
		}()
	}

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(saveCh)

	}()

	// Read from saveCh and process responses
	for resp := range saveCh {
		saver.Save(resp)
	}

	// Close all connections in the pool
	close(connPool)
	for conn := range connPool {
		conn.Disconnect()
	}
}
