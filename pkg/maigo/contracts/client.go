package contracts

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/jeanmolossi/maigo/pkg/maigo/header"
	"github.com/jeanmolossi/maigo/pkg/maigo/mime"
)

// Client combines configuration and the ability to create HTTP requests for
// different methods. Implementations manage shared resources like headers or
// cookies and expose fluent builders for each HTTP verb.
type Client interface {
	ClientConfig
	ClientHTTPMethods
}

// ClientCompat exposes the MaiGo client configuration and provides access to
// the underlying *http.Client for interoperability with the Go standard
// library.
type ClientCompat interface {
	Client
	// Unwrap exposes the configured *http.Client, carrying over MaiGo's
	// client-level settings such as timeout, transport and redirect
	// behaviour.
	Unwrap() *http.Client
}

// ClientBuilder builds a Client. It allows configuring default headers and
// cookies before producing the final ClientHTTPMethods value.
//
// Example:
//
//	client := maigo.NewClient("https://api.example.com").
//	        Header().Add(header.Accept, mime.JSON.String()).
//	        Build()
type ClientBuilder interface {
	// Header returns a builder to configure default headers.
	Header() BuilderHeader[ClientBuilder]
	// Cookie returns a builder to configure default cookies.
	Cookie() BuilderCookie[ClientBuilder]
	// Build finalizes the configuration and produces a ClientHTTPMethods.
	Build() ClientHTTPMethods
}

// ClientConfig exposes the configurable parts of a client, such as the
// underlying HTTP client, default headers, cookies, validations and base URL.
type ClientConfig interface {
	ConfigHTTPClient
	// Header exposes the client's default headers.
	Header() Header
	// Cookies exposes the client's cookie jar.
	Cookies() Cookies
	// Validations returns collected validation errors.
	Validations() Validations
	ConfigBaseURL
}

// ConfigHTTPClient allows replacing or retrieving the underlying HTTP client.
type ConfigHTTPClient interface {
	// SetHttpClient replaces the underlying HTTP client implementation.
	SetHttpClient(httpc HTTPClient)
	// HttpClient retrieves the current HTTP client.
	HttpClient() HTTPClient
}

// HTTPClientCompat is the minimal interface implemented by http.Client. It
// adds an Unwrap method to access the concrete *http.Client value configured by
// MaiGo.
type HTTPClientCompat interface {
	HTTPClient
	// Unwrap exposes the configured *http.Client, carrying over MaiGo's
	// client-level settings such as timeout, transport and redirect
	// behaviour.
	Unwrap() *http.Client
}

// HTTPClient is the minimal interface implemented by http.Client. It enables
// callers to plug in custom clients with specific transports or timeouts.
type HTTPClient interface {
	// Do sends an HTTP request and returns an HTTP response.
	Do(r *http.Request) (*http.Response, error)
	// Transport returns the RoundTripper used by the client.
	Transport() http.RoundTripper
	// SetTransport sets the RoundTripper used by the client.
	SetTransport(rt http.RoundTripper)
	// Timeout returns the request timeout.
	Timeout() time.Duration
	// SetTimeout sets the request timeout.
	SetTimeout(d time.Duration)
	// SetFollowRedirects controls whether redirects are automatically followed.
	SetFollowRedirects(follow bool)
}

// ClientHTTPMethods defines fluent builders for each HTTP method. Each method
// returns a RequestBuilder prepared with the respective verb and path.
type ClientHTTPMethods interface {
	// GET prepares a GET request for the given path.
	GET(path string) RequestBuilder
	// POST prepares a POST request for the given path.
	POST(path string) RequestBuilder
	// PUT prepares a PUT request for the given path.
	PUT(path string) RequestBuilder
	// DELETE prepares a DELETE request for the given path.
	DELETE(path string) RequestBuilder
	// PATCH prepares a PATCH request for the given path.
	PATCH(path string) RequestBuilder
	// HEAD prepares a HEAD request for the given path.
	HEAD(path string) RequestBuilder
	// CONNECT prepares a CONNECT request for the given path.
	CONNECT(path string) RequestBuilder
	// OPTIONS prepares an OPTIONS request for the given path.
	OPTIONS(path string) RequestBuilder
	// TRACE prepares a TRACE request for the given path.
	TRACE(path string) RequestBuilder
}

// ConfigBaseURL exposes the base URL used to resolve request paths.
type ConfigBaseURL interface {
	// BaseURL returns the base URL used to resolve request paths.
	BaseURL() *url.URL
}

// BuilderHeader configures HTTP headers for the parent builder. Each method
// returns the parent type so calls can be chained.
//
// Example:
//
//	builder.Header().
//	        Add(header.Accept, mime.JSON.String()).
//	        Set(header.UserAgent, "mai-go")
type BuilderHeader[T any] interface {
	// Add appends a header value.
	Add(key header.Type, value string) T
	// AddAll appends all provided headers.
	AddAll(headers map[header.Type]string) T
	// Set replaces the header value.
	Set(key header.Type, value string) T
	// SetAll replaces all provided headers.
	SetAll(headers map[header.Type]string) T
	// AddAccept adds an Accept header with the given mime type.
	AddAccept(value mime.Type) T
	// AddContentType adds a Content-Type header with the given mime type.
	AddContentType(value mime.Type) T
	// AddUserAgent adds a User-Agent header.
	AddUserAgent(value string) T
}

// BuilderCookie adds cookies to the parent builder.
type BuilderCookie[T any] interface {
	// Add includes a cookie to be sent with requests.
	Add(cookie *http.Cookie) T
}

// BuilderAuth configures authentication headers like Bearer or Basic auth.
type BuilderAuth[T any] interface {
	// Set writes the Authorization header as provided.
	Set(value string) T
	// BearerToken sets the Authorization header using a bearer token.
	BearerToken(token string) T
	// BasicAuth sets the Authorization header using basic auth credentials.
	BasicAuth(user, pass string) T
}

// BuilderHTTPClientConfig tunes the behaviour of the underlying HTTP client,
// allowing custom transports, timeouts and proxy settings.
type BuilderHTTPClientConfig[T any] interface {
	// SetCustomHTTPClient swaps the default HTTP client implementation.
	SetCustomHTTPClient(httpClient HTTPClient) T
	// SetCustomTransport overrides the client's transport.
	SetCustomTransport(transport http.RoundTripper) T
	// SetTimeout defines the request timeout.
	SetTimeout(duration time.Duration) T
	// SetFollowRedirects determines whether redirects are followed automatically.
	SetFollowRedirects(follow bool) T
	// SetProxy configures an HTTP proxy via URL.
	SetProxy(proxyURL string) T
}

// BuilderRequestContext sets the context used when sending a request.
type BuilderRequestContext[T any] interface {
	// Set defines the context to use when the request is sent.
	Set(ctx context.Context) T
}

// BuilderRequestBody serializes values into the request body.
// Supported formats include raw readers, strings, JSON and XML.
type BuilderRequestBody[T any] interface {
	// AsReader uses the raw reader as the request body.
	AsReader(body io.Reader) T
	// AsString writes the provided string as the request body.
	AsString(body string) T
	// AsJSON serializes obj as JSON into the request body.
	AsJSON(obj any) T
	// AsXML serializes obj as XML into the request body.
	AsXML(obj any) T
}
