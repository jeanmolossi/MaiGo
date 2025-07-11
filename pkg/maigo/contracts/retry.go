package contracts

import (
	"time"
)

type BuilderRequestRetry[T any] interface {
	SetConstantBackoff(interval time.Duration, maxAttempts uint) T
	SetConstantBackoffWithJitter(interval time.Duration, maxAttempts uint) T
	SetExponentialBackoff(interval time.Duration, maxAttempts uint, backoffRate float64) T
	SetExponentialBackoffWithJitter(interval time.Duration, maxAttempts uint, backoffRate float64) T
	WithRetryCondition(shouldRetry func(response Response) bool) T
	WithMaxDelay(duration time.Duration) T
}
