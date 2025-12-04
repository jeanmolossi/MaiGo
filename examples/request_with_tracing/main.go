// Package main demonstrates MaiGo client tracing with OpenTelemetry.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jeanmolossi/maigo/examples/testserver"
	"github.com/jeanmolossi/maigo/pkg/httpx"
	"github.com/jeanmolossi/maigo/pkg/httpx/tracing"
	"github.com/jeanmolossi/maigo/pkg/maigo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	stdoutmetric "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.opentelemetry.io/otel/trace"
)

func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := stdout.New(
		stdout.WithPrettyPrint(),
		stdout.WithoutTimestamps(),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp, nil
}

func initMeter() (*sdkmetric.MeterProvider, error) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")

	exp, err := stdoutmetric.New(
		stdoutmetric.WithEncoder(enc),
		stdoutmetric.WithoutTimestamps(),
	)
	if err != nil {
		return nil, err
	}

	res := resource.NewSchemaless(semconv.ServiceName("stdoutmetric-maigo-example"))

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp)),
		sdkmetric.WithResource(res),
	)

	otel.SetMeterProvider(mp)

	return mp, nil
}

func main() {
	tp, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}

	mp, err := initMeter()
	if err != nil {
		log.Fatal(err)
	}

	ts := testserver.NewManager().NewServer()
	defer ts.Close()

	bag, _ := baggage.Parse("username=johndoe")
	ctx := baggage.ContextWithBaggage(context.Background(), bag)

	var body []byte

	tr := otel.Tracer("examples/request_with_tracing")

	transport := httpx.Compose(
		http.DefaultTransport,
		tracing.WithTracing(),
	)

	client := maigo.NewClient(ts.URL).
		Config().SetCustomTransport(transport).
		Build()

	err = func(ctx context.Context) error {
		ctx, span := tr.Start(ctx, "say hello", trace.WithAttributes(semconv.PeerService("ExampleService")))
		defer span.End()

		span.SetAttributes(semconv.URLFull(ts.URL))
		span.SetAttributes(semconv.HTTPRequestMethodKey.String(http.MethodGet))

		res, err := client.GET("/").Context().Set(ctx).Send()
		if err != nil {
			return err
		}

		body, err = res.Body().AsBytes()

		return err
	}(ctx)
	if err != nil {
		log.Fatal(err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = tp.Shutdown(shutdownCtx)
	if err != nil {
		panic(err)
	}

	err = mp.Shutdown(shutdownCtx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("response received: %s\n\n\n", body)
}
