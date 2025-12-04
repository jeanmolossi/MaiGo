package tracing

import (
	"fmt"
	"net/http"

	"github.com/jeanmolossi/maigo/pkg/httpx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.opentelemetry.io/otel/trace"
)

const tracerName = "github.com/jeanmolossi/maigo/pkg/httpx/tracing"

// WithTracing starts a span for each outbound HTTP request and propagates the
// span context via request headers.
func WithTracing() httpx.ChainedRoundTripper {
	tracer := otel.Tracer(tracerName)
	propagator := otel.GetTextMapPropagator()

	return func(next http.RoundTripper) http.RoundTripper {
		return httpx.RoundTripperFn(func(req *http.Request) (*http.Response, error) {
			ctx := req.Context()

			ctx, span := tracer.Start(
				ctx,
				fmt.Sprintf("%s %s", req.Method, req.URL.RequestURI()),
				trace.WithSpanKind(trace.SpanKindClient),
			)
			defer span.End()

			req = req.Clone(ctx)
			propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))

			resp, err := next.RoundTrip(req)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())

				return nil, err
			}

			span.SetAttributes(
				semconv.HTTPRequestMethodKey.String(req.Method),
				semconv.URLFullKey.String(req.URL.String()),
				semconv.HTTPResponseStatusCodeKey.Int(resp.StatusCode),
				attribute.String("http.target", req.URL.Path),
			)

			return resp, nil
		})
	}
}
