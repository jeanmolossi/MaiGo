package testserver

import "time"

type Server struct {
	ID                int
	URL               string
	EnableBusy        bool
	EnableHeaderDebug bool

	Interval    time.Duration
	BackoffRate float64
}
