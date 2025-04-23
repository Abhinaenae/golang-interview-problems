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

	saveCh := make(chan string, len(requests))
	var wg sync.WaitGroup

	for _, req := range requests {
		conn, err := creator.NewConnection()
		if err != nil {
			slog.Error("Failed to create connection", "error", err)
			continue
		}
		defer conn.Disconnect()
		conn.Connect()

		wg.Add(1)
		go func(r string) {
			defer wg.Done()
			resp, err := conn.Send(r)
			if err != nil {
				slog.Error("Failed to send request", "error", err)
				return
			}
			saveCh <- resp
		}(req)
	}

	wg.Wait()
	close(saveCh)

	for resp := range saveCh {
		saver.Save(resp)
	}

}
