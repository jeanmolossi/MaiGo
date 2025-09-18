package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"

	"github.com/jeanmolossi/maigo/examples/testserver"
	"github.com/jeanmolossi/maigo/pkg/httpx"
	"github.com/jeanmolossi/maigo/pkg/httpx/metrics"
	"github.com/jeanmolossi/maigo/pkg/maigo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	ts := testserver.NewManager().NewServer()
	defer ts.Close()

	registry := prometheus.NewRegistry()

	metricsChain := metrics.MetricsRoundTripper(metrics.RoundTripperOptions{
		Registerer: registry,
		Namespace:  "maigo",
		Subsystem:  "client",
	})

	transport := httpx.Compose(http.DefaultTransport, metricsChain)

	client := maigo.NewClient(ts.URL).
		Config().
		SetCustomTransport(transport).
		Build()

	_, err := client.GET("/users").Send()
	if err != nil {
		slog.Error("request failed", "error", err)
	}

	metricsHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	recorder := httptest.NewRecorder()
	metricsHandler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/metrics", nil))

	fmt.Print(recorder.Body.String())
}
