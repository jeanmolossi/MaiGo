package maigo

import (
	"time"

	"github.com/jeanmolossi/maigo/pkg/maigo/contracts"
)

var _ contracts.BuilderRequestRetry[contracts.RequestBuilder] = (*RequestRetryBuilder)(nil)

type RequestRetryBuilder struct {
	parent        *RequestBuilder
	requestConfig *RequestConfigBase
}

func (r *RequestBuilder) Retry() contracts.BuilderRequestRetry[contracts.RequestBuilder] {
	return &RequestRetryBuilder{
		parent:        r,
		requestConfig: r.request.config,
	}
}

// SetConstantBackoff implements contracts.BuilderRequestRetry.
func (r *RequestRetryBuilder) SetConstantBackoff(interval time.Duration, maxAttempts uint) contracts.RequestBuilder {
	r.requestConfig.RetryConfig().SetInterval(interval)
	r.requestConfig.RetryConfig().SetMaxAttempts(maxAttempts)
	r.requestConfig.RetryConfig().SetBackoffRate(1)
	r.requestConfig.RetryConfig().SetJitterStrategy(JitterStrategyNone)

	return r.parent
}

// SetConstantBackoffWithJitter implements contracts.BuilderRequestRetry.
func (r *RequestRetryBuilder) SetConstantBackoffWithJitter(interval time.Duration, maxAttempts uint) contracts.RequestBuilder {
	r.requestConfig.RetryConfig().SetInterval(interval)
	r.requestConfig.RetryConfig().SetMaxAttempts(maxAttempts)
	r.requestConfig.RetryConfig().SetBackoffRate(1)
	r.requestConfig.RetryConfig().SetJitterStrategy(JitterStrategyFull)

	return r.parent
}

// SetExponentialBackoff implements contracts.BuilderRequestRetry.
func (r *RequestRetryBuilder) SetExponentialBackoff(interval time.Duration, maxAttempts uint, backoffRate float64) contracts.RequestBuilder {
	r.requestConfig.RetryConfig().SetInterval(interval)
	r.requestConfig.RetryConfig().SetMaxAttempts(maxAttempts)
	r.requestConfig.RetryConfig().SetBackoffRate(backoffRate)
	r.requestConfig.RetryConfig().SetJitterStrategy(JitterStrategyNone)

	return r.parent
}

// SetExponentialBackoffWithJitter implements contracts.BuilderRequestRetry.
func (r *RequestRetryBuilder) SetExponentialBackoffWithJitter(interval time.Duration, maxAttempts uint, backoffRate float64) contracts.RequestBuilder {
	r.requestConfig.RetryConfig().SetInterval(interval)
	r.requestConfig.RetryConfig().SetMaxAttempts(maxAttempts)
	r.requestConfig.RetryConfig().SetBackoffRate(backoffRate)
	r.requestConfig.RetryConfig().SetJitterStrategy(JitterStrategyFull)

	return r.parent
}

// WithMaxDelay implements contracts.BuilderRequestRetry.
func (r *RequestRetryBuilder) WithMaxDelay(duration time.Duration) contracts.RequestBuilder {
	r.requestConfig.RetryConfig().SetMaxDelay(duration)
	return r.parent
}

// WithRetryCondition implements contracts.BuilderRequestRetry.
func (r *RequestRetryBuilder) WithRetryCondition(shouldRetry func(response contracts.Response) bool) contracts.RequestBuilder {
	r.requestConfig.RetryConfig().SetShouldRetry(shouldRetry)
	return r.parent
}
