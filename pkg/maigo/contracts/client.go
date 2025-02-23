package contracts

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/jeanmolossi/MaiGo/pkg/maigo/header"
	"github.com/jeanmolossi/MaiGo/pkg/maigo/mime"
)

type Client interface {
	ClientConfig
	ClientHTTPMethods
}

type ClientBuilder interface {
	Header() BuilderHeader[ClientBuilder]
	Cookie() BuilderCookie[ClientBuilder]
	Build() ClientHTTPMethods
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

type BuilderHeader[T any] interface {
	Add(key header.Type, value string) T
	AddAll(headers map[header.Type]string) T
	Set(key header.Type, value string) T
	SetAll(headers map[header.Type]string) T
	AddAccept(value mime.Type) T
	AddContentType(value mime.Type) T
	AddUserAgent(value string) T
}

type BuilderCookie[T any] interface {
	Add(cookie *http.Cookie) T
}

type BuilderAuth[T any] interface {
	Set(value string) T
	BearerToken(token string) T
	BasicAuth(user, pass string) T
}

type BuilderHTTPClientConfig[T any] interface {
	SetCustomHTTPClient(httpClient HTTPClient) T
	SetCustomTransport(transport http.RoundTripper) T
	SetTimeout(duration time.Duration) T
	SetFollowRedirects(follow bool) T
	SetProxy(proxyURL string) T
}

type BuilderRequestContext[T any] interface {
	Set(ctx context.Context) T
}

type BuilderRequestBody[T any] interface {
	AsReader(body io.Reader) T
	AsString(body string) T
	AsJSON(obj any) T
	AsXML(obj any) T
}
