package retry

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jeanmolossi/MaiGo/pkg/httpx"
	"github.com/stretchr/testify/require"
)

func backoffZero(int) time.Duration { return 0 }

// ----- Tests -----

// Replays small body faithfully and sets attempt header (using POST to avoid the AllowedMethods bug).
func TestRetry_ReplaysSmallBody_AndSetsAttemptHeader(t *testing.T) {
	base, assert := httpx.NewRoundTripMockBuilder().
		AddOutcome(httpx.NewResp(500, "fail"), nil).
		AddOutcome(httpx.NewResp(200, "ok"), nil).
		Build(t)

	cfg := RetryConfig{
		MaxAttempts:        3,
		Backoff:            backoffZero,
		IgnoreRetryAfter:   true,
		MaxReplayBodyBytes: 64 << 10,
		ReplayBodyStrategy: ReplayIfSmallElseNoRetry,
	}

	rt := WithRetry(cfg)(base)

	req, _ := http.NewRequest(http.MethodPost, "http://x", bytes.NewBufferString("hello"))
	req.Header = make(http.Header)
	req = req.WithContext(context.Background())

	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp == nil || resp.StatusCode != 200 {
		t.Fatalf("expected final 200, got %#v", resp)
	}

	assert.Calls(2, "expected 2 calls (500 then 200)")
	assert.SeenBodiesLen(2, "expected to record bodies for 2 attempts")
	assert.SeenBodies(0, "hello", "body replay mismatch")
	assert.SeenBodies(1, "hello", "body replay mismatch")
	assert.SeenHeaders(0, defaultAttemptHeader, "1")
	assert.SeenHeaders(1, defaultAttemptHeader, "2")
}

// Respect Retry-After with cap (use POST para entrar no fluxo de retry).
func TestRetry_RespectsRetryAfter_WithCap(t *testing.T) {
	var observedDelay []time.Duration

	base, assert := httpx.NewRoundTripMockBuilder().
		AddOutcome(
			httpx.NewResponseBuilder(503, "").
				SetHeader("Retry-After", "120").
				Build(),
			nil,
		).
		AddOutcome(httpx.NewResp(200, ""), nil).
		Build(t)

	cfg := RetryConfig{
		MaxAttempts:        2,
		Backoff:            backoffZero,
		IgnoreRetryAfter:   false,
		MaxRetryAfter:      50 * time.Millisecond, // cap
		ReplayBodyStrategy: ReplayIfSmallElseNoRetry,
	}
	cfg.OnRetry = func(ctx context.Context, attempt int, r *http.Request, resp *http.Response, err error, delay time.Duration) {
		observedDelay = append(observedDelay, delay)
	}

	rt := WithRetry(cfg)(base)

	req, _ := http.NewRequest(http.MethodPost, "http://x", bytes.NewBufferString("a"))

	_, err := rt.RoundTrip(req)
	require.NoError(t, err, "unexpected: %v", err)
	require.Len(t, observedDelay, 1, "expected exactly one retry delay observed, got %d", len(observedDelay))
	require.Equal(t, 50*time.Millisecond, observedDelay[0], "expected delay caped to 50ms, got %v", observedDelay[0])
	assert.Calls(2)
}

// Ignore Retry-After (delay must come from Backoff).
func TestRetry_IgnoreRetryAfter(t *testing.T) {
	var observedDelay []time.Duration

	base, assert := httpx.NewRoundTripMockBuilder().
		AddOutcome(
			httpx.NewResponseBuilder(503, "").
				SetHeader("Retry-After", "120").
				Build(),
			nil,
		).
		AddOutcome(httpx.NewResp(200, ""), nil).
		Build(t)

	cfg := RetryConfig{
		MaxAttempts: 2,
		Backoff: func(attempt int) time.Duration {
			return 250 * time.Millisecond
		},
		IgnoreRetryAfter:   true,
		MaxRetryAfter:      5 * time.Second,
		ReplayBodyStrategy: ReplayIfSmallElseNoRetry,
	}
	cfg.OnRetry = func(ctx context.Context, attempt int, r *http.Request, resp *http.Response, err error, delay time.Duration) {
		observedDelay = append(observedDelay, delay)
	}

	rt := WithRetry(cfg)(base)

	req, _ := http.NewRequest(http.MethodPost, "http://x", bytes.NewBufferString("a"))

	_, err := rt.RoundTrip(req)
	require.NoError(t, err, "unexpected: %v", err)
	require.Len(t, observedDelay, 1, "expected delay from backoff")
	require.Equal(t, 250*time.Millisecond, observedDelay[0], "expected delay from backoff (250ms), got %v (%d obs)", observedDelay, len(observedDelay))
	assert.Calls(2)
}

type noGetBodyReader struct{ r io.Reader }

func (n noGetBodyReader) Read(p []byte) (int, error) { return n.r.Read(p) }

// When body is bigger than MaxReplayBodyBytes and strategy does NOT spill, do not retry (only one attempt).
func TestRetry_BodyTooBig_NoRetry_WhenNoSpill(t *testing.T) {
	// 70KiB body > 64KiB default
	huge := bytes.Repeat([]byte("x"), 70<<10)

	base, assert := httpx.NewRoundTripMockBuilder().
		AddOutcome(httpx.NewResp(500, ""), nil).
		AddOutcome(httpx.NewResp(200, ""), nil).
		Build(t)

	cfg := RetryConfig{
		MaxAttempts:        3,
		Backoff:            backoffZero,
		IgnoreRetryAfter:   true,
		MaxReplayBodyBytes: 64 << 10,                 // 64KiB
		ReplayBodyStrategy: ReplayIfSmallElseNoRetry, // sem spill
	}

	rt := WithRetry(cfg)(base)

	body := io.NopCloser(noGetBodyReader{r: bytes.NewReader(huge)})

	req, _ := http.NewRequest(http.MethodPut, "http://x", body)
	resp, err := rt.RoundTrip(req)
	require.NoError(t, err, "unexpected error: %v", err)
	require.Equal(t, 500, resp.StatusCode, "expected response 500, got %d", resp.StatusCode)
	assert.Calls(1, "expected 1 call (should not retry with large body and no spill)")
}

func TestRetry_AllowedMethods_RetryOnGET(t *testing.T) {
	base, assert := httpx.NewRoundTripMockBuilder().
		AddOutcome(httpx.NewResp(500, ""), nil).
		AddOutcome(httpx.NewResp(200, ""), nil).
		Build(t)

	cfg := RetryConfig{
		MaxAttempts:        2,
		Backoff:            backoffZero,
		IgnoreRetryAfter:   true,
		ReplayBodyStrategy: ReplayIfSmallElseNoRetry,
		// defaultAllowed inclui GET
	}

	rt := WithRetry(cfg)(base)

	req, _ := http.NewRequest(http.MethodGet, "http://x", nil)
	resp, err := rt.RoundTrip(req)
	require.NoError(t, err, "unexpected: %v", err)
	require.Equal(t, 200, resp.StatusCode, "expected to return last response (200) with retry, got %d", resp.StatusCode)
	assert.Calls(2)
}

// Sanity: Attempt header is decimal "1", "2", ... and resets per request.
func TestRetry_AttemptHeaderSequence(t *testing.T) {
	base, assert := httpx.NewRoundTripMockBuilder().
		AddOutcome(httpx.NewResp(500, ""), nil).
		AddOutcome(httpx.NewResp(200, ""), nil).
		Build(t)

	cfg := RetryConfig{
		MaxAttempts:        2,
		Backoff:            backoffZero,
		IgnoreRetryAfter:   true,
		ReplayBodyStrategy: ReplayIfSmallElseNoRetry,
	}

	rt := WithRetry(cfg)(base)

	req, _ := http.NewRequest(http.MethodPost, "http://x", strings.NewReader("k"))
	_, _ = rt.RoundTrip(req)

	assert.SeenHeadersLen(2, "expected two attempts with headers recorded")
	assert.SeenHeaders(0, defaultAttemptHeader, "1", "attempt header value wanted 1")
	assert.SeenHeaders(1, defaultAttemptHeader, "2", "attempt header value wanted 2")
}

// parseRetryAfter sanity (seconds and HTTP-date). Uses OnRetry to read computed delay.
func TestRetry_parseRetryAfter_SecondsAndDate(t *testing.T) {
	// 1) seconds form
	{
		var delays []time.Duration

		base, assert := httpx.NewRoundTripMockBuilder().
			AddOutcome(
				httpx.NewResponseBuilder(503, "").
					SetHeader("Retry-After", "1").
					Build(),
				nil,
			).
			AddOutcome(httpx.NewResp(200, ""), nil).
			Build(t)

		cfg := RetryConfig{
			MaxAttempts:        2,
			Backoff:            backoffZero,
			IgnoreRetryAfter:   false,
			MaxRetryAfter:      10 * time.Second,
			ReplayBodyStrategy: ReplayIfSmallElseNoRetry,
		}
		cfg.OnRetry = func(_ context.Context, _ int, _ *http.Request, _ *http.Response, _ error, d time.Duration) {
			delays = append(delays, d)
		}
		rt := WithRetry(cfg)(base)

		req, _ := http.NewRequest(http.MethodPost, "http://x", strings.NewReader("a"))
		_, _ = rt.RoundTrip(req)

		require.Len(t, delays, 1, "Retry-After seconds not respected: %#v", delays)
		require.Equal(t, 1*time.Second, delays[0])
		assert.Calls(2)
	}

	// 2) HTTP-date form (now+2s → expect cap/min handled)
	{
		when := time.Now().Add(1 * time.Second).UTC().Format(http.TimeFormat)

		var delays []time.Duration

		base, assert := httpx.NewRoundTripMockBuilder().
			AddOutcome(
				httpx.NewResponseBuilder(503, "").
					SetHeader("Retry-After", when).
					Build(),
				nil,
			).
			AddOutcome(httpx.NewResp(200, ""), nil).
			Build(t)

		cfg := RetryConfig{
			MaxAttempts:        2,
			Backoff:            backoffZero,
			IgnoreRetryAfter:   false,
			MaxRetryAfter:      10 * time.Second,
			ReplayBodyStrategy: ReplayIfSmallElseNoRetry,
			OnRetry: func(ctx context.Context, attempt int, r *http.Request, resp *http.Response, err error, d time.Duration) {
				delays = append(delays, d)
			},
		}

		rt := WithRetry(cfg)(base)

		req, _ := http.NewRequest(http.MethodPost, "http://x", strings.NewReader("a"))
		_, _ = rt.RoundTrip(req)

		require.Len(t, delays, 1)

		// Aceita +/- 1s por execução em tempo real
		if delays[0] < 50*time.Millisecond || delays[0] > 1250*time.Millisecond {
			t.Fatalf("unexpected HTTP-date delay: %v", delays[0])
		}

		assert.Calls(2)
	}
}

// ----- Bench-lite sanity (optional): verify no sleep with backoffZero -----

func TestRetry_NoSleep_WhenBackoffZero(t *testing.T) {
	base, assert := httpx.NewRoundTripMockBuilder().
		AddOutcome(httpx.NewResp(500, ""), nil).
		AddOutcome(httpx.NewResp(200, ""), nil).
		Build(t)

	cfg := RetryConfig{
		MaxAttempts:        2,
		Backoff:            backoffZero,
		IgnoreRetryAfter:   true,
		ReplayBodyStrategy: ReplayIfSmallElseNoRetry,
	}
	start := time.Now()

	rt := WithRetry(cfg)(base)

	req, _ := http.NewRequest(http.MethodPut, "http://x", strings.NewReader("x"))
	_, _ = rt.RoundTrip(req)

	require.Less(t, time.Since(start), 1*time.Millisecond, "unexpected delay without backoff/Retry-After")
	assert.Calls(2)
}
