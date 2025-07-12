package testserver

import "time"

type Server struct {
	ID                int
	URL               string
	EnableBusy        bool
	EnableHeaderDebug bool

	// Interval specifies the base sleep duration for request delays
	Interval time.Duration
	// BackoffRate specifies the multiplier for exponential backoff (e.g., 2.0 for doubling delay)
	BackoffRate float64
}
