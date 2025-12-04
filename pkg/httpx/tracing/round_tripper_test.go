package tracing

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/jeanmolossi/maigo/pkg/httpx"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func TestWithTracingStartsAndEndsSpan(t *testing.T) {
	spanRecorder, restore := installTracer(t)
	t.Cleanup(restore)

	req := mustRequest(http.MethodGet, "https://example.com/test")
	res := httpx.NewResponseBuilder(http.StatusOK, "").SetRequest(req).Build()

	mock, assert := httpx.NewRoundTripMockBuilder().AddOutcome(res, nil).Build(t)
	transport := httpx.Compose(mock, WithTracing())

	resp, err := transport.RoundTrip(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Calls(1)

	require.Len(t, spanRecorder.Ended(), 1)
	span := spanRecorder.Ended()[0]
	require.Equal(t, codes.Unset, span.Status().Code)
	require.Equal(t, trace.SpanKindClient, span.SpanKind())
}

func TestWithTracingPropagatesContextAndRecordsErrors(t *testing.T) {
	spanRecorder, restore := installTracer(t)
	t.Cleanup(restore)

	roundTripErr := errors.New("transport failure")
	mock, assert := httpx.NewRoundTripMockBuilder().AddOutcome(nil, roundTripErr).Build(t)
	transport := httpx.Compose(mock, WithTracing())

	ctx := context.Background()
	tracer := otel.Tracer("test/tracer")

	ctx, parent := tracer.Start(ctx, "parent")
	defer parent.End()

	req := mustRequest(http.MethodPost, "https://example.com/propagate")
	req = req.WithContext(ctx)

	_, err := transport.RoundTrip(req)
	require.Error(t, err)
	assert.Calls(1)
	assert.SeenHeadersLen(1)

	ended := spanRecorder.Ended()
	require.Len(t, ended, 1)
	span := ended[0]

	carrier := propagation.MapCarrier{}
	propagation.TraceContext{}.Inject(trace.ContextWithSpanContext(context.Background(), span.SpanContext()), carrier)

	assert.SeenHeaders(0, "traceparent", carrier.Get("traceparent"))
	require.Equal(t, codes.Error, span.Status().Code)
	require.True(t, span.SpanContext().TraceID().IsValid())
}

func installTracer(t *testing.T) (*tracetest.SpanRecorder, func()) {
	t.Helper()

	spanRecorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(spanRecorder),
	)

	prevProvider := otel.GetTracerProvider()
	prevPropagator := otel.GetTextMapPropagator()

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	restore := func() {
		_ = tp.Shutdown(context.Background())

		otel.SetTracerProvider(prevProvider)
		otel.SetTextMapPropagator(prevPropagator)
	}

	return spanRecorder, restore
}

func mustRequest(method, rawURL string) *http.Request {
	req, err := http.NewRequest(method, rawURL, http.NoBody)
	if err != nil {
		panic(err)
	}

	return req
}
