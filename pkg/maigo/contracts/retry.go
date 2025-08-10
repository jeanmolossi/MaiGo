package contracts

import (
	"time"
)

// BuilderRequestRetry configures retry logic for a request. It supports
// constant or exponential backoff strategies, optional jitter, custom retry
// conditions and maximum delay between attempts.
//
// Example:
//
//	builder.Retry().
//	        SetExponentialBackoff(500*time.Millisecond, 3, 2.0)
type BuilderRequestRetry[T any] interface {
	// SetConstantBackoff retries at a fixed interval for the given attempts.
	SetConstantBackoff(interval time.Duration, maxAttempts uint) T
	// SetConstantBackoffWithJitter retries at a fixed interval plus jitter.
	SetConstantBackoffWithJitter(interval time.Duration, maxAttempts uint) T
	// SetExponentialBackoff retries with exponentially increasing intervals.
	SetExponentialBackoff(interval time.Duration, maxAttempts uint, backoffRate float64) T
	// SetExponentialBackoffWithJitter retries with exponential backoff and jitter.
	SetExponentialBackoffWithJitter(interval time.Duration, maxAttempts uint, backoffRate float64) T
	// WithRetryCondition retries only when shouldRetry returns true.
	WithRetryCondition(shouldRetry func(response Response) bool) T
	// WithMaxDelay caps the total retry delay.
	WithMaxDelay(duration time.Duration) T
}
