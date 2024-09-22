package main

import (
	"time"
)

// FlipFlopTimer represents a timer that alternates between on and off states
type FlipFlopTimer struct {
	On         bool
	Duration   time.Duration
	LastChange time.Time
}

// NewFlipFlopTimer creates a new FlipFlopTimer with the given duration
func NewFlipFlopTimer(duration time.Duration) *FlipFlopTimer {
	return &FlipFlopTimer{
		On:         false,
		Duration:   duration,
		LastChange: time.Now(),
	}
}

// Tick updates the state of the timer based on the elapsed time
func (f *FlipFlopTimer) Tick() {
	if time.Since(f.LastChange) > f.Duration {
		f.LastChange = time.Now()
		f.On = !f.On
	}
}

func main() {
	// Example usage
	timer := NewFlipFlopTimer(2 * time.Second)

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		timer.Tick()
		if timer.On {
			println("Timer is ON")
		} else {
			println("Timer is OFF")
		}
	}
}
