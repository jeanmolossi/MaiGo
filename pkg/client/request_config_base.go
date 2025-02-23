package client

import (
	"net/url"
	"time"

	"github.com/jeanmolossi/MaiGo/pkg/client/contracts"
	"github.com/jeanmolossi/MaiGo/pkg/client/method"
)

type (
	RequestConfigBase struct {
		ctx          contracts.Context
		httpHeader   contracts.Header
		httpCookies  contracts.Cookies
		method       method.Type
		path         string
		searchParams url.Values
		body         any
		validations  contracts.Validations
		retryConfig  *RetryConfig
	}

	JitterStrategy string

	RetryConfig struct {
		shouldRetry    func(response contracts.Response) bool
		interval       time.Duration
		maxAttempts    uint
		backoffRate    float64
		maxDelay       *time.Duration
		jitterStrategy JitterStrategy
	}
)

const (
	// JitterStrategyNone NONE is the default jitter strategy.
	JitterStrategyNone JitterStrategy = "NONE"
	// JitterStrategyFull FULL is the full jitter strategy.
	JitterStrategyFull JitterStrategy = "FULL"
)

func (r *RequestConfigBase) Context() contracts.Context {
	return r.ctx
}

func (r *RequestConfigBase) Header() contracts.Header {
	return r.httpHeader
}

func (r *RequestConfigBase) Method() method.Type {
	return r.method
}

func (r *RequestConfigBase) Path() string {
	return r.path
}

func (r *RequestConfigBase) SearchParams() url.Values {
	return r.searchParams
}

func (r *RequestConfigBase) RetryConfig() *RetryConfig {
	return r.retryConfig
}

// RetryConfig methods

func (r *RetryConfig) ShouldRetry() func(response contracts.Response) bool {
	return r.shouldRetry
}

func (r *RetryConfig) SetShouldRetry(shouldRetry func(response contracts.Response) bool) {
	r.shouldRetry = shouldRetry
}

func (r *RetryConfig) Interval() time.Duration {
	return r.interval
}

func (r *RetryConfig) SetInterval(duration time.Duration) {
	r.interval = duration
}

func (r *RetryConfig) MaxAttempts() uint {
	return r.maxAttempts
}

func (r *RetryConfig) SetMaxAttempts(attempts uint) {
	r.maxAttempts = attempts
}

func (r *RetryConfig) BackoffRate() float64 {
	return r.backoffRate
}

func (r *RetryConfig) SetBackoffRate(rate float64) {
	r.backoffRate = rate
}

func (r *RetryConfig) MaxDelay() *time.Duration {
	return r.maxDelay
}

func (r *RetryConfig) SetMaxDelay(duration time.Duration) {
	r.maxDelay = &duration
}

func (r *RetryConfig) JitterStrategy() JitterStrategy {
	return r.jitterStrategy
}

func (r *RetryConfig) SetJitterStrategy(strategy JitterStrategy) {
	r.jitterStrategy = strategy
}

func newRequestConfigBase(method method.Type, path string) *RequestConfigBase {
	return &RequestConfigBase{
		ctx:          newDefaultContext(),
		httpHeader:   newDefaultHTTPHeader(),
		httpCookies:  newDefaultHttpCookies(),
		method:       method,
		path:         path,
		searchParams: url.Values{},
		body:         newBufferedBody(),
		validations:  newDefaultValidations(nil),
		retryConfig: &RetryConfig{
			shouldRetry: func(response contracts.Response) bool {
				return response.Status().IsError()
			},
			interval:       1 * time.Second,
			backoffRate:    2.0,
			jitterStrategy: JitterStrategyNone,
		},
	}
}
