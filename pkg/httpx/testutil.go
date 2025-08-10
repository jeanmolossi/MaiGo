package httpx

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type outcome struct {
	resp *http.Response
	err  error
}

type roundTripMock struct {
	mu sync.Mutex

	outcomes []outcome

	calls       int
	seenBodies  []string
	seenHeaders []http.Header
}

func (s *roundTripMock) RoundTrip(r *http.Request) (*http.Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var bodyBytes []byte
	if r.Body != nil && r.Body != http.NoBody {
		bodyBytes, _ = io.ReadAll(r.Body)
		_ = r.Body.Close()
	}

	s.seenBodies = append(s.seenBodies, string(bodyBytes))
	s.seenHeaders = append(s.seenHeaders, r.Header.Clone())

	if s.calls >= len(s.outcomes) {
		s.calls++

		resp := NewResp(200, "")
		resp.Request = r

		return resp, nil
	}

	out := s.outcomes[s.calls]
	s.calls++

	return out.resp, out.err
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

func (r *RoundTripMockBuilder) Build(t *testing.T) (http.RoundTripper, Assert) {
	if t == nil {
		return r.mock, Assert{}
	}

	t.Helper()

	return r.mock, Assert{mock: r.mock, t: t}
}

// Assert

type Assert struct {
	mock *roundTripMock
	t    *testing.T
}

func (a Assert) Calls(expected int, args ...any) {
	require.Equal(a.t, expected, a.mock.calls, args...)
}

func (a Assert) SeenBodiesLen(expected int, args ...any) {
	require.Len(a.t, a.mock.seenBodies, expected, args...)
}

func (a Assert) SeenBodies(callIdx int, expected string, args ...any) {
	require.Equal(a.t, expected, a.mock.seenBodies[callIdx], args...)
}

func (a Assert) SeenHeadersLen(expected int, args ...any) {
	require.Len(a.t, a.mock.seenHeaders, expected, args...)
}

func (a Assert) SeenHeaders(callIdx int, key, expected string, args ...any) {
	require.Equal(a.t, expected, a.mock.seenHeaders[callIdx].Get(key), args...)
}

// ResponseBuilder

type ResponseBuilder struct {
	resp *http.Response
}

func NewResponseBuilder(status int, body string) *ResponseBuilder {
	resp := NewResp(status, body)

	return &ResponseBuilder{
		resp: resp,
	}
}

func (r *ResponseBuilder) SetRequest(req *http.Request) *ResponseBuilder {
	r.resp.Request = req
	return r
}

func (r *ResponseBuilder) SetHeader(k, v string) *ResponseBuilder {
	r.resp.Header.Set(k, v)
	return r
}

func (r *ResponseBuilder) Build() *http.Response {
	return r.resp
}

func NewResp(status int, body string) *http.Response {
	headers := make(http.Header)

	if len(body) > 0 {
		headers.Set("Content-Length", strconv.Itoa(len(body)))
	}

	return &http.Response{
		StatusCode: status,
		Status:     strconv.Itoa(status) + " " + http.StatusText(status),
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     headers,
	}
}
