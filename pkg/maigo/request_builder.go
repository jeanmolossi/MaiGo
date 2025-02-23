package maigo

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"

	"github.com/jeanmolossi/MaiGo/pkg/maigo/contracts"
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

	for attempt := range config.MaxAttempts() {
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
		time.Sleep(delay)
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

func secureFloat64() float64 {
	var bits [8]byte
	// read randomic 8 bytes from secure source
	_, err := rand.Read(bits[:])
	if err != nil {
		panic(err)
	}

	// parse bytes to 64 bits int
	n := binary.LittleEndian.Uint64(bits[:])

	// normalize to interval [0, 1]
	//nolint:mnd // 64 shift
	return float64(n) / (1 << 64)
}
