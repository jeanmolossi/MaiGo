package circuitbreaker

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/jeanmolossi/MaiGo/pkg/httpx"
	"github.com/stretchr/testify/require"
)

func TestCircuitBreaker_OpenAndClose(t *testing.T) {
	base, assert := httpx.NewRoundTripMockBuilder().
		AddOutcome(nil, errors.New("fail")).
		AddOutcome(nil, errors.New("fail")).
		AddOutcome(nil, errors.New("fail")).
		AddOutcome(httpx.NewResp(200, "ok"), nil).
		AddOutcome(httpx.NewResp(200, "ok"), nil).
		Build(t)

	cfg := CircuitBreakerConfig{FailureThreshold: 3, RecoveryWindow: 50 * time.Millisecond}
	rt := WithCircuitBreaker(cfg)(base)

	req, _ := http.NewRequest(http.MethodGet, "http://x", nil)

	for i := 0; i < 3; i++ {
		_, err := rt.RoundTrip(req)
		require.Error(t, err)
	}

	_, err := rt.RoundTrip(req)
	require.ErrorIs(t, err, ErrCircuitOpen)
	assert.Calls(3)

	time.Sleep(60 * time.Millisecond)

	resp, err := rt.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	assert.Calls(4)

	resp, err = rt.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	assert.Calls(5)
}

func TestCircuitBreaker_HalfOpenFailureReopens(t *testing.T) {
	base, assert := httpx.NewRoundTripMockBuilder().
		AddOutcome(nil, errors.New("fail")).
		AddOutcome(nil, errors.New("fail")).
		AddOutcome(nil, errors.New("fail")).
		AddOutcome(nil, errors.New("fail")).
		Build(t)

	cfg := CircuitBreakerConfig{FailureThreshold: 3, RecoveryWindow: 50 * time.Millisecond}
	rt := WithCircuitBreaker(cfg)(base)

	req, _ := http.NewRequest(http.MethodGet, "http://x", nil)

	for i := 0; i < 3; i++ {
		_, err := rt.RoundTrip(req)
		require.Error(t, err)
	}

	_, err := rt.RoundTrip(req)
	require.ErrorIs(t, err, ErrCircuitOpen)
	assert.Calls(3)

	time.Sleep(60 * time.Millisecond)

	_, err = rt.RoundTrip(req)
	require.Error(t, err)
	assert.Calls(4)

	_, err = rt.RoundTrip(req)
	require.ErrorIs(t, err, ErrCircuitOpen)
	assert.Calls(4)
}
