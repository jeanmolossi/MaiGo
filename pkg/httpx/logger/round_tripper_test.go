package logger_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jeanmolossi/MaiGo/pkg/httpx"
	"github.com/jeanmolossi/MaiGo/pkg/httpx/logger"
	"github.com/stretchr/testify/require"
)

type memLogger struct {
	info []string
	errs []string
}

func (m *memLogger) Info(ctx context.Context, msg string, args ...any) { m.info = append(m.info, msg) }
func (m *memLogger) Error(ctx context.Context, err error, msg string, args ...any) {
	m.errs = append(m.errs, err.Error())
}

const mockURL = "http://foo.bar/baz"

func TestLogger_HooksAndBodySample(t *testing.T) {
	t.Parallel()

	next, _ := httpx.NewRoundTripMockBuilder().
		AddOutcome(httpx.NewResp(200, `{"pong":true}`), nil).
		Build(t)

	log := &memLogger{}

	var started, ended bool

	rt := httpx.Compose(next,
		logger.LoggerRoundTripper(logger.LoggerHooks{
			Logger:  log,
			OnStart: func(ctx context.Context, r *http.Request) { started = true },
			OnEnd:   func(ctx context.Context, r *http.Request, resp *http.Response, err error) { ended = true },
			ReqBodyTransformerFn: func(ctx context.Context) func([]byte) []byte {
				return func(b []byte) []byte { return bytes.ToUpper(b) }
			},
			ResBodyTransformerFn: func(ctx context.Context) func([]byte) []byte {
				return func(b []byte) []byte { return b }
			},
			LogStart: true,
			LogEnd:   true,
		}),
	)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, mockURL, bytes.NewBufferString(`{"ping":true}`))
	require.NoError(t, err, "unexpected error: %s", err)

	resp, err := rt.RoundTrip(req)
	require.NoError(t, err, "unexpected error: %s", err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "expected %s, got %s", "OK", http.StatusText(resp.StatusCode))
	require.True(t, started, "expected request started")
	require.True(t, ended, "expected request ended")
	require.Len(t, log.info, 2, "expected 2 log infos, got %d", len(log.info))
	require.Empty(t, log.errs, "expected log errors empty, got %d", len(log.errs))
}

func TestLogger_DefaultHooks(t *testing.T) {
	t.Parallel()

	next, _ := httpx.NewRoundTripMockBuilder().
		AddOutcome(httpx.NewResp(200, `{"pong":true}`), nil).
		Build(t)

	log := &memLogger{}

	rt := httpx.Compose(next,
		logger.LoggerRoundTripper(logger.LoggerHooks{
			Logger:   log,
			LogStart: true,
			LogEnd:   true,
		}),
	)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, mockURL, bytes.NewBufferString(`{"ping":true}`))
	require.NoError(t, err, "unexpected error: %s", err)

	resp, err := rt.RoundTrip(req)
	require.NoError(t, err, "unexpected error: %s", err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "expected %s, got %s", "OK", http.StatusText(resp.StatusCode))
	require.Len(t, log.info, 2, "expected 2 log infos, got %d", len(log.info))
	require.Empty(t, log.errs, "expected log errors empty, got %d", len(log.errs))
}

func TestLogger_Errors(t *testing.T) {
	t.Parallel()

	t.Run("when fail request, should have error log", testFailingRequest)
	t.Run("when done request without content-length header, should have log error", testDoneRequestWithoutHeaders)
}

func testFailingRequest(t *testing.T) {
	t.Parallel()

	next, _ := httpx.NewRoundTripMockBuilder().
		AddOutcome(
			httpx.NewResponseBuilder(500, `{}`).
				SetHeader("Content-Length", "0"). // proposital content different from body
				Build(),
			errors.New("request error"),
		).
		Build(t)

	log := &memLogger{}

	rt := httpx.Compose(next,
		logger.LoggerRoundTripper(logger.LoggerHooks{
			Logger:   log,
			LogStart: true,
			LogEnd:   true,
		}),
	)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, mockURL, bytes.NewBufferString(`{"ping":true}`))
	require.NoError(t, err, "unexpected error: %s", err)

	resp, err := rt.RoundTrip(req)
	require.NotNil(t, resp, "unexpected nil response")
	require.EqualError(t, err, "request error")
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode, "expected %s, got %s", "Internal Server Error", http.StatusText(resp.StatusCode))
	require.Len(t, log.info, 1, "expected 1 log infos, got %d", len(log.info))
	require.Len(t, log.errs, 1, "expected 1 errors")
}

func testDoneRequestWithoutHeaders(t *testing.T) {
	t.Parallel()

	mockResp := httpx.NewResp(200, `{}`)
	mockResp.Header = make(http.Header)

	next, _ := httpx.NewRoundTripMockBuilder().
		AddOutcome(mockResp, nil).
		Build(t)

	log := &memLogger{}

	rt := httpx.Compose(next,
		logger.LoggerRoundTripper(logger.LoggerHooks{
			Logger:   log,
			LogStart: true,
			LogEnd:   true,
		}),
	)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, mockURL, bytes.NewBufferString(`{"ping":true}`))
	require.NoError(t, err, "unexpected error: %s", err)

	resp, err := rt.RoundTrip(req)
	require.NoError(t, err, "unexpected error: %s", err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "expected %s, got %s", "OK", http.StatusText(resp.StatusCode))
	require.Len(t, log.info, 2, "expected 2 log infos, got %d", len(log.info))
	require.Len(t, log.errs, 1, "expected empty errors")
}

func TestLoggerRoundTripper_TruncatesRequestBodyAndLogsSize(t *testing.T) {
	const maxLogBodySize = 65536 // 64KiB

	// Server reads ALL who client sents and reports how much was received.
	var received int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close() //nolint:errcheck

		b, _ := io.ReadAll(r.Body)
		received = len(b)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	t.Cleanup(srv.Close)

	// Payload extrapolate max to enforce truncating.
	payload := bytes.Repeat([]byte("A"), maxLogBodySize+12345)

	var read int

	// Transformers who NOT delete body, else you log empty.
	hooks := logger.LoggerHooks{
		LogStart:       true,
		LogEnd:         true,
		SupressErrors:  true,
		MaxLogBodySize: maxLogBodySize,
		ReqBodyTransformerFn: func(ctx context.Context) func([]byte) []byte {
			return func(b []byte) []byte { read = len(b); return b }
		},
		ResBodyTransformerFn: func(ctx context.Context) func([]byte) []byte {
			return func(b []byte) []byte { return b }
		},
		StartMessage: "http.client.start",
		EndMessage:   "http.client.end",
	}

	// Chain our LogRoundTripper with default transporter
	rt := logger.LoggerRoundTripper(hooks)(http.DefaultTransport)

	req, err := http.NewRequest(http.MethodPost, srv.URL, bytes.NewReader(payload))
	require.NoError(t, err)

	resp, err := rt.RoundTrip(req)
	require.NoError(t, err)

	_ = resp.Body.Close()

	sent := maxLogBodySize + 12345
	require.Equal(t, sent, received)       // ensure the server received full body
	require.Equal(t, maxLogBodySize, read) // ensure the log transformer received body truncated
}
