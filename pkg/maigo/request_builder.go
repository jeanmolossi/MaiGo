package maigo

import (
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	mrand "math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/jeanmolossi/MaiGo/pkg/maigo/contracts"
	"github.com/jeanmolossi/MaiGo/pkg/maigo/header"
)

var _ contracts.RequestBuilder = (*RequestBuilder)(nil)

var (
	secureRand   = newSecureRand()
	secureRandMu sync.Mutex
)

type RequestBuilder struct {
	request *Request
}

func (r *RequestBuilder) createFullURL() *url.URL {
	// parse base URL and path
	fullURL := r.request.client.BaseURL().JoinPath(r.request.config.Path())

	query := fullURL.Query()

	for param, values := range r.request.config.SearchParams() {
		for _, value := range values {
			query.Add(param, value)
		}
	}

	fullURL.RawQuery = query.Encode()

	return fullURL
}

func (r *RequestBuilder) createHTTPRequest() (*http.Request, error) {
	// create full URL
	fullURL := r.createFullURL()

	request, err := http.NewRequestWithContext(
		r.request.config.Context().Unwrap(),
		r.request.config.Method().String(),
		fullURL.String(),
		r.request.config.body.Unwrap(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	// Add client cookies
	for _, cookie := range r.request.client.Cookies().Unwrap() {
		request.AddCookie(cookie)
	}

	// Add request cookies
	for _, cookie := range r.request.config.Cookies().Unwrap() {
		request.AddCookie(cookie)
	}

	// Add client headers
	for key, values := range *r.request.client.Header().Unwrap() {
		for _, value := range values {
			request.Header.Set(key, value)
		}
	}

	// Add request headers
	for key, values := range *r.request.config.Header().Unwrap() {
		for _, value := range values {
			request.Header.Set(key, value)
		}
	}

	return request, nil
}

func (r *RequestBuilder) execute(request *http.Request) (contracts.Response, error) {
	//nolint:bodyclose // newResponse method reads response.Body, then it can not be closed here
	response, err := r.request.client.HttpClient().Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return newResponse(response), nil
}

func (r *RequestBuilder) executeWithRetry(request *http.Request) (contracts.Response, error) {
	config := r.request.config.RetryConfig()

	//nolint:prealloc
	var (
		executionErr error
		attemptsErr  []error
		response     contracts.Response
	)

	const retryHeader = header.Type("X-Retry-Attempt")

	for attempt := range config.MaxAttempts() {
		request.Header.Set(retryHeader.String(), strconv.FormatUint(uint64(attempt), 36))

		response, executionErr = r.execute(request)

		if executionErr == nil {
			shouldRetry := config.ShouldRetry()

			if !shouldRetry(response) {
				return response, nil
			}

			//nolint:err113
			executionErr = errors.New(response.Status().Text())
		}

		attemptsErr = append(attemptsErr, fmt.Errorf("[call %d]: %w", attempt+1, executionErr))

		// delay before another try
		delay := r.calculateRetryDelay(attempt)
		if err := sleepCtx(request.Context(), delay); err != nil {
			return nil, err
		}
	}

	return nil, fmt.Errorf(
		"request failed after %d attempts: %w",
		config.MaxAttempts(),
		errors.Join(attemptsErr...),
	)
}

func (r *RequestBuilder) calculateRetryDelay(attempt uint) time.Duration {
	config := r.request.config.RetryConfig()

	delay := float64(config.Interval()) * math.Pow(config.BackoffRate(), float64(attempt))

	if config.MaxDelay() != nil {
		delay = math.Min(delay, float64(*config.MaxDelay()))
	}

	if config.JitterStrategy() == JitterStrategyFull {
		delay = secureFloat64() * delay
	}

	return time.Duration(delay)
}

func (r *RequestBuilder) Send() (contracts.Response, error) {
	if err := errors.Join(r.request.client.Validations().Unwrap()...); err != nil {
		return nil, errors.Join(ErrClientValidation, err)
	}

	if err := errors.Join(r.request.config.Validations().Unwrap()...); err != nil {
		return nil, errors.Join(ErrRequestValidation, err)
	}

	req, err := r.createHTTPRequest()
	if err != nil {
		return nil, errors.Join(ErrCreateRequest, err)
	}

	retry := r.request.config.RetryConfig()
	if retry != nil && retry.MaxAttempts() > 1 {
		return r.executeWithRetry(req)
	}

	return r.execute(req)
}

func newSecureRand() *mrand.Rand {
	var b [8]byte
	if _, err := crand.Read(b[:]); err != nil {
		panic(err)
	}

	seed := int64(binary.LittleEndian.Uint64(b[:]))

	return mrand.New(mrand.NewSource(seed))
}

func secureFloat64() float64 {
	secureRandMu.Lock()
	v := secureRand.Float64()
	secureRandMu.Unlock()

	return v
}

func sleepCtx(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}

	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
