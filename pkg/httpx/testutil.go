package httpx

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
)

type outcome struct {
	resp *http.Response
	err  error
}

type roundTripMock struct {
	outcomes []outcome
	idx      int32
}

func (s *roundTripMock) RoundTrip(r *http.Request) (*http.Response, error) {
	i := atomic.AddInt32(&s.idx, 1) - 1
	if int(i) >= len(s.outcomes) {
		i = int32(len(s.outcomes) - 1)
	}

	o := s.outcomes[i]

	return o.resp, o.err
}

type RoundTripMockBuilder struct {
	mock *roundTripMock
}

func NewRoundTripMockBuilder() *RoundTripMockBuilder {
	return &RoundTripMockBuilder{
		&roundTripMock{},
	}
}

func (r *RoundTripMockBuilder) AddOutcome(res *http.Response, err error) *RoundTripMockBuilder {
	r.mock.outcomes = append(r.mock.outcomes, outcome{resp: res, err: err})
	return r
}

func (r *RoundTripMockBuilder) Build() http.RoundTripper {
	return r.mock
}

func NewResp(status int, body string) *http.Response {
	headers := make(http.Header)

	if len(body) > 0 {
		headers.Set("Content-Length", strconv.Itoa(len(body)))
	}

	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     headers,
	}
}
