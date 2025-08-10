package retry

import (
	"context"
	"errors"
	"io"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jeanmolossi/MaiGo/pkg/httpx"
)

const (
	defaultAttemptHeader = "X-Retry-Attempt"
	defaultInterval      = 100 * time.Millisecond
	defaultBackoffRate   = 2
	maxBackoff           = float64(5 * time.Second)
	maxReplayBodyBytes   = 64 << 10 // 64KiB
)

type RetryConfig struct {
	MaxAttempts    int
	AllowedMethods map[string]bool
	ShouldRetry    func(*http.Request, *http.Response, error) bool
	Backoff        func(attempt int) time.Duration
	OnRetry        func(ctx context.Context, attempt int, r *http.Request, resp *http.Response, err error, delay time.Duration)

	IgnoreRetryAfter bool
	MaxRetryAfter    time.Duration
	AttemptHeader    string

	MaxReplayBodyBytes int
	ReplayBodyStrategy BodyReplayStrategy
}

func WithRetry(cfg RetryConfig) httpx.ChainedRoundTripper {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 3
	}

	if cfg.AttemptHeader == "" {
		cfg.AttemptHeader = defaultAttemptHeader
	}

	if cfg.AllowedMethods == nil {
		cfg.AllowedMethods = defaultAllowed()
	}

	if cfg.ShouldRetry == nil {
		cfg.ShouldRetry = defaultShouldRetry
	}

	if cfg.Backoff == nil {
		cfg.Backoff = defaultBackoff
	}

	if cfg.MaxReplayBodyBytes <= 0 {
		cfg.MaxReplayBodyBytes = maxReplayBodyBytes
	}

	if cfg.ReplayBodyStrategy == 0 {
		cfg.ReplayBodyStrategy = ReplayIfSmallElseNoRetry
	}

	if cfg.MaxRetryAfter == 0 {
		cfg.MaxRetryAfter = 30 * time.Second
	}

	return func(next http.RoundTripper) http.RoundTripper {
		return httpx.RoundTripperFn(func(r *http.Request) (*http.Response, error) {
			if allowed, ok := cfg.AllowedMethods[strings.ToUpper(r.Method)]; !ok || !allowed {
				return next.RoundTrip(r)
			}

			req := httpx.CloneRequest(r)

			cleanup, bodyOK, berr := ensureReopenableBody(req, int64(cfg.MaxReplayBodyBytes), cfg.ReplayBodyStrategy)
			if berr != nil {
				return next.RoundTrip(r) // do not retry and use original request
			}

			if cleanup != nil {
				defer cleanup()
			}

			allowRetryWithBody := bodyOK || !hasRequestBody(req)

			var (
				resp *http.Response
				err  error
			)

			for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
				// renew context
				req = req.WithContext(r.Context())

				// retry attempt header
				req.Header.Set(cfg.AttemptHeader, strconv.FormatUint(uint64(attempt), 10))

				// reopen body (if have)
				if req.GetBody != nil {
					nb, gerr := req.GetBody()
					if gerr != nil {
						return nil, gerr
					}

					req.Body = nb
				}

				resp, err = next.RoundTrip(req)

				if !cfg.ShouldRetry(req, resp, err) || attempt == cfg.MaxAttempts || !allowRetryWithBody {
					return resp, err
				}

				delay := cfg.Backoff(attempt)

				if !cfg.IgnoreRetryAfter &&
					err == nil &&
					resp != nil &&
					(resp.StatusCode == 429 || resp.StatusCode == 503) {
					if ra, ok := parseRetryAfter(resp.Header.Get("Retry-After")); ok {
						delay = min(ra, cfg.MaxRetryAfter)
					}
				}

				if cfg.OnRetry != nil {
					cfg.OnRetry(req.Context(), attempt, req, resp, err, delay)
				}

				if resp != nil && resp.Body != nil {
					_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20)) // drains until 1MiB
					_ = resp.Body.Close()
				}

				if serr := sleepCtx(req.Context(), delay); serr != nil {
					return resp, serr
				}
			}

			// unreachable code
			return resp, err
		})
	}
}

func hasRequestBody(r *http.Request) bool {
	return r.Body != nil && r.Body != http.NoBody
}

func defaultAllowed() map[string]bool {
	return map[string]bool{
		http.MethodGet:     true,
		http.MethodHead:    true,
		http.MethodPut:     true,
		http.MethodDelete:  true,
		http.MethodOptions: true,
		http.MethodTrace:   true,
	}
}

func defaultShouldRetry(_ *http.Request, resp *http.Response, err error) bool {
	if err != nil {
		var ne net.Error

		if errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
			return true
		}

		if errors.As(err, &ne) && ne.Timeout() {
			return true
		}

		msg := strings.ToLower(err.Error())
		if strings.Contains(msg, "connection reset") || strings.Contains(msg, "broken pipe") ||
			strings.Contains(msg, "refused") || strings.Contains(msg, "unexpected eof") {
			return true
		}

		return true // ??
	}

	if resp == nil {
		return true
	}

	switch resp.StatusCode {
	case 408, 425, 429, 500, 502, 503, 504:
		return true
	default:
		return false
	}
}

func defaultBackoff(attempt int) time.Duration {
	delay := float64(defaultInterval) * math.Pow(defaultBackoffRate, float64(attempt))
	delay = math.Min(delay, maxBackoff)

	return time.Duration(delay)
}

func sleepCtx(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}

	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
