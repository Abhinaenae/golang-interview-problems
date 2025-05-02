package main

import (
	"sync/atomic"
	"time"
)

type Monitor interface {
	SendRate(int)
}

type Handler struct {
	count atomic.Uint64
}

func (h *Handler) Handle() {
	// Some work
	// Calculate rate of this function
	h.count.Add(1)

}

// LogRate starts a goroutine that periodically reports the number of handled events
// to the provided monitor. It samples the count every duration d and resets the counter
// after each report. The monitoring continues until the program terminates.
func (h *Handler) LogRate(monitor Monitor, d time.Duration) {
	go func() {
		ticker := time.NewTicker(d)
		defer ticker.Stop()

		for range ticker.C {
			currCount := h.count.Swap(0)
			monitor.SendRate(int(currCount))

		}
	}()
}
