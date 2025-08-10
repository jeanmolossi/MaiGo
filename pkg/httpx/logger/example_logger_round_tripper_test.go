package logger_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/jeanmolossi/MaiGo/pkg/httpx"
	"github.com/jeanmolossi/MaiGo/pkg/httpx/logger"
)

// ExampleLoggerRoundTripper demonstrates how to chain the logging round tripper
// with a transport. It logs request/response metadata and returns the mocked
// response from the builder.
func ExampleLoggerRoundTripper() {
	base, _ := httpx.NewRoundTripMockBuilder().
		AddOutcome(httpx.NewResp(200, `{"pong":true}`), nil).
		Build(new(testing.T))

	rt := httpx.Compose(
		base,
		logger.LoggerRoundTripper(logger.LoggerHooks{
			Logger:   logger.NewNoop(),
			LogStart: true,
			LogEnd:   true,
		}),
	)

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "http://foo.bar", bytes.NewBufferString(`{"ping":true}`))
	resp, _ := rt.RoundTrip(req)
	fmt.Println(resp.StatusCode)
	// Output:
	// 200
}
