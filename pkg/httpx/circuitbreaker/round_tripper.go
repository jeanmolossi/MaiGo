package circuitbreaker

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/jeanmolossi/MaiGo/pkg/httpx"
)

// ErrCircuitOpen is returned when the circuit is open and requests are not allowed.
var ErrCircuitOpen = errors.New("circuit breaker is open")

type state int

const (
	stateClosed state = iota
	stateOpen
	stateHalfOpen
)

// CircuitBreakerConfig holds settings for the circuit breaker middleware.
type CircuitBreakerConfig struct {
	// FailureThreshold defines how many consecutive failures are allowed before
	// the circuit opens.
	FailureThreshold int
	// RecoveryWindow is the time the circuit remains open before allowing a
	// probe request in half-open state.
	RecoveryWindow time.Duration
	// ShouldTrip determines whether a response/error should be considered a
	// failure. If nil, errors or >=500 responses trip the circuit.
	ShouldTrip func(*http.Response, error) bool
}

// WithCircuitBreaker wraps the next RoundTripper with circuit breaker logic.
func WithCircuitBreaker(cfg CircuitBreakerConfig) httpx.ChainedRoundTripper {
	if cfg.FailureThreshold <= 0 {
		cfg.FailureThreshold = 5
	}

	if cfg.RecoveryWindow <= 0 {
		cfg.RecoveryWindow = 30 * time.Second
	}

	if cfg.ShouldTrip == nil {
		cfg.ShouldTrip = defaultShouldTrip
	}

	return func(next http.RoundTripper) http.RoundTripper {
		return &cbTransport{next: next, cfg: cfg}
	}
}

type cbTransport struct {
	next http.RoundTripper
	cfg  CircuitBreakerConfig

	mu       sync.Mutex
	state    state
	failures int
	openedAt time.Time
	probing  bool
}

func (c *cbTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	c.mu.Lock()
	switch c.state {
	case stateOpen:
		if time.Since(c.openedAt) >= c.cfg.RecoveryWindow {
			c.state = stateHalfOpen
		} else {
			c.mu.Unlock()
			return nil, ErrCircuitOpen
		}
	}

	if c.state == stateHalfOpen {
		if c.probing {
			c.mu.Unlock()
			return nil, ErrCircuitOpen
		}

		c.probing = true
	}
	c.mu.Unlock()

	resp, err := c.next.RoundTrip(r)

	trip := c.cfg.ShouldTrip(resp, err)

	c.mu.Lock()
	defer c.mu.Unlock()

	switch c.state {
	case stateHalfOpen:
		c.probing = false
		if trip {
			c.state = stateOpen
			c.openedAt = time.Now()
			c.failures = 0
		} else {
			c.state = stateClosed
			c.failures = 0
		}
	case stateClosed:
		if trip {
			c.failures++
			if c.failures >= c.cfg.FailureThreshold {
				c.state = stateOpen
				c.openedAt = time.Now()
				c.failures = 0
			}
		} else {
			c.failures = 0
		}
	}

	return resp, err
}

func defaultShouldTrip(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}

	if resp == nil {
		return true
	}

	return resp.StatusCode >= 500
}
