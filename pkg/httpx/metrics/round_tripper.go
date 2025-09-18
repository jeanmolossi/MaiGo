package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/jeanmolossi/maigo/pkg/httpx"
	"github.com/prometheus/client_golang/prometheus"
)

// RoundTripperOptions configures the MetricsRoundTripper behaviour.
type RoundTripperOptions struct {
	// Registerer is used to register the collectors. Defaults to prometheus.DefaultRegisterer.
	Registerer prometheus.Registerer

	// DurationBuckets allows overriding the buckets used by the duration histogram.
	DurationBuckets []float64

	// Namespace is prefixed to the metric names.
	Namespace string
	// Subsystem is added to the metric names after the namespace.
	Subsystem string

	// DurationName overrides the default duration metric name.
	DurationName string
	// CountName overrides the default request counter name.
	CountName string

	// DurationCollector allows providing a pre-constructed histogram vector.
	DurationCollector *prometheus.HistogramVec
	// CountCollector allows providing a pre-constructed counter vector.
	CountCollector *prometheus.CounterVec
}

const (
	defaultDurationName = "request_duration_seconds"
	defaultCountName    = "requests_total"
)

// MetricsRoundTripper instruments an HTTP client transport recording request
// durations and counts labelled by method and status code.
func MetricsRoundTripper(opts RoundTripperOptions) httpx.ChainedRoundTripper {
	duration := opts.DurationCollector
	count := opts.CountCollector

	registerer := opts.Registerer
	if registerer == nil {
		registerer = prometheus.DefaultRegisterer
	}

	if duration == nil {
		duration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      metricName(opts.DurationName, defaultDurationName),
			Help:      "Duration of outbound HTTP requests",
			Buckets:   bucketsOrDefault(opts.DurationBuckets),
		}, []string{"method", "status"})
	}

	duration = registerHistogram(registerer, duration)

	if count == nil {
		count = prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      metricName(opts.CountName, defaultCountName),
			Help:      "Total number of outbound HTTP requests",
		}, []string{"method", "status"})
	}

	count = registerCounter(registerer, count)

	return func(next http.RoundTripper) http.RoundTripper {
		return httpx.RoundTripperFn(func(r *http.Request) (*http.Response, error) {
			start := time.Now()
			resp, err := next.RoundTrip(r)
			elapsed := time.Since(start).Seconds()

			status := "error"
			if err == nil && resp != nil {
				status = strconv.Itoa(resp.StatusCode)
			}

			labels := prometheus.Labels{
				"method": r.Method,
				"status": status,
			}

			duration.With(labels).Observe(elapsed)
			count.With(labels).Inc()

			return resp, err
		})
	}
}

func metricName(provided, fallback string) string {
	if provided != "" {
		return provided
	}

	return fallback
}

func bucketsOrDefault(buckets []float64) []float64 {
	if len(buckets) > 0 {
		return buckets
	}

	return prometheus.DefBuckets
}

func registerHistogram(registerer prometheus.Registerer, collector *prometheus.HistogramVec) *prometheus.HistogramVec {
	if registerer == nil {
		return collector
	}

	if err := registerer.Register(collector); err != nil {
		if alreadyRegistered, ok := err.(prometheus.AlreadyRegisteredError); ok {
			if existing, ok := alreadyRegistered.ExistingCollector.(*prometheus.HistogramVec); ok {
				return existing
			}

			return collector
		}
	}

	return collector
}

func registerCounter(registerer prometheus.Registerer, collector *prometheus.CounterVec) *prometheus.CounterVec {
	if registerer == nil {
		return collector
	}

	if err := registerer.Register(collector); err != nil {
		if alreadyRegistered, ok := err.(prometheus.AlreadyRegisteredError); ok {
			if existing, ok := alreadyRegistered.ExistingCollector.(*prometheus.CounterVec); ok {
				return existing
			}

			return collector
		}
	}

	return collector
}
