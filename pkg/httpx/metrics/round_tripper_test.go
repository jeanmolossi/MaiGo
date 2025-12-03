package metrics

import (
	"errors"
	"net/http"
	"testing"

	"github.com/jeanmolossi/maigo/pkg/httpx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/require"
)

func TestMetricsRoundTripper_IncrementsSuccessMetrics(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	counter := prometheus.NewCounterVec(prometheus.CounterOpts{Name: "maigo_test_requests_total"}, []string{"method", "status"})
	duration := prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "maigo_test_request_duration_seconds"}, []string{"method", "status"})

	rt := MetricsRoundTripper(RoundTripperOptions{
		Registerer:        registry,
		CountCollector:    counter,
		DurationCollector: duration,
	})

	next := httpx.RoundTripperFn(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusCreated, Body: http.NoBody}, nil
	})

	req, err := http.NewRequest(http.MethodPost, "http://example.com", http.NoBody)
	require.NoError(t, err)

	_, err = rt(next).RoundTrip(req)
	require.NoError(t, err)

	require.InEpsilon(t, 1, testutil.ToFloat64(counter.WithLabelValues(http.MethodPost, "201")), 0.0001)
	requireHistogramSampleCount(t, registry, "maigo_test_request_duration_seconds", http.MethodPost, "201", 1)
}

func TestMetricsRoundTripper_IncrementsErrorMetrics(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	counter := prometheus.NewCounterVec(prometheus.CounterOpts{Name: "maigo_test_requests_total"}, []string{"method", "status"})
	duration := prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "maigo_test_request_duration_seconds"}, []string{"method", "status"})

	rt := MetricsRoundTripper(RoundTripperOptions{
		Registerer:        registry,
		CountCollector:    counter,
		DurationCollector: duration,
	})

	next := httpx.RoundTripperFn(func(r *http.Request) (*http.Response, error) {
		return nil, errors.New("boom")
	})

	req, err := http.NewRequest(http.MethodGet, "http://example.com", http.NoBody)
	require.NoError(t, err)

	_, err = rt(next).RoundTrip(req)
	require.Error(t, err)

	require.InEpsilon(t, 1, testutil.ToFloat64(counter.WithLabelValues(http.MethodGet, "error")), 0.0001)
	requireHistogramSampleCount(t, registry, "maigo_test_request_duration_seconds", http.MethodGet, "error", 1)
}

func requireHistogramSampleCount(t *testing.T, registry *prometheus.Registry, metricName, method, status string, expected uint64) {
	t.Helper()

	metricFamilies, err := registry.Gather()
	require.NoError(t, err)

	for _, mf := range metricFamilies {
		if mf.GetName() != metricName {
			continue
		}

		for _, m := range mf.GetMetric() {
			if labelsMatch(m, method, status) {
				require.NotNil(t, m.GetHistogram(), "expected histogram for metric %s", metricName)
				require.Equal(t, expected, m.GetHistogram().GetSampleCount())
				require.Greater(t, m.GetHistogram().GetSampleSum(), 0.0)

				return
			}
		}
	}

	t.Fatalf("metric %s with method=%s status=%s not found", metricName, method, status)
}

func labelsMatch(m *dto.Metric, method, status string) bool {
	var foundMethod, foundStatus bool

	for _, lp := range m.GetLabel() {
		switch lp.GetName() {
		case "method":
			if lp.GetValue() == method {
				foundMethod = true
			}
		case "status":
			if lp.GetValue() == status {
				foundStatus = true
			}
		}
	}

	return foundMethod && foundStatus
}
