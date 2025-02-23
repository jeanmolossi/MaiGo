package contracts

import (
	"net/http"
	"net/url"
	"time"
)

type Client interface {
	ClientConfig
	ClientHTTPMethods
}

type ClientConfig interface {
	ConfigHTTPClient
	Header() Header
	Cookies() Cookies
	Validations() Validations
	ConfigBaseURL
}

type ConfigHTTPClient interface {
	SetHttpClient(httpc HTTPClient)
	HttpClient() HTTPClient
}

type HTTPClient interface {
	Do(r *http.Request) (*http.Response, error)
	Transport() http.RoundTripper
	SetTransport(rt http.RoundTripper)
	Timeout() time.Duration
	SetTimeout(d time.Duration)
	SetFollowRedirects(follow bool)
}

type ClientHTTPMethods interface {
	GET(path string) RequestBuilder
	POST(path string) RequestBuilder
	PUT(path string) RequestBuilder
	DELETE(path string) RequestBuilder
	PATCH(path string) RequestBuilder
	HEAD(path string) RequestBuilder
	CONNECT(path string) RequestBuilder
	OPTIONS(path string) RequestBuilder
	TRACE(path string) RequestBuilder
}

type ConfigBaseURL interface {
	BaseURL() *url.URL
}
