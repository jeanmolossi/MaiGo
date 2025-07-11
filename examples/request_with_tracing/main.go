package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jeanmolossi/MaiGo/examples/testserver"
	"github.com/jeanmolossi/MaiGo/pkg/maigo"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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
	exporter, err := stdout.New(stdout.WithPrettyPrint())
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

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	ts := testserver.NewManager().NewServer()
	defer ts.Close()

	bag, _ := baggage.Parse("username=johndoe")
	ctx := baggage.ContextWithBaggage(context.Background(), bag)

	var body []byte

	tr := otel.Tracer("examples/request_with_tracing")

	oteltransporter := otelhttp.NewTransport(nil, otelhttp.WithMeterProvider(mp))

	client := maigo.NewClient(ts.URL).
		Config().SetCustomTransport(oteltransporter).
		Build()

	err = func(ctx context.Context) error {
		ctx, span := tr.Start(ctx, "say hello", trace.WithAttributes(semconv.PeerService("ExampleService")))
		defer span.End()

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

	err = tp.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}

	err = mp.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("response received: %s\n\n\n", body)
}
