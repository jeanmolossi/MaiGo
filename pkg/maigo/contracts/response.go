package contracts

import (
	"io"
	"net/http"
)

// Response represents an HTTP response and provides fluent helpers for
// inspecting its body, headers, cookies, originating request and status.
type Response interface {
	// Raw exposes the underlying http.Response.
	Raw() *http.Response
	// Body provides helpers for reading the response body.
	Body() ResponseFluentBody
	// Header provides helpers for inspecting response headers.
	Header() ResponseFluentHeader
	// Cookie provides helpers for accessing response cookies.
	Cookie() ResponseFluentCookie
	// Request describes the request that generated this response.
	Request() ResponseFluentRequest
	// Status provides helpers for checking the HTTP status code.
	Status() ResponseFluentStatus
}

// ResponseFluentBody exposes helpers to read the response body in various
// formats. Calling any of the read methods consumes and closes the underlying
// body.
//
// Example:
//
//	text, _ := resp.Body().AsString()
type ResponseFluentBody interface {
	// Raw returns the underlying body reader.
	Raw() io.Closer
	// Close closes the body reader.
	Close()
	// AsBytes reads the body as raw bytes.
	AsBytes() ([]byte, error)
	// AsString reads the body as a string.
	AsString() (string, error)
        // AsJSON decodes the body as JSON into v.
        AsJSON(v any) error
}

// ResponseFluentCookie provides access to cookies returned by the server.
type ResponseFluentCookie interface {
	// GetAll returns all cookies returned by the server.
	GetAll() []*http.Cookie
}

// ResponseFluentHeader allows querying HTTP headers from the response.
type ResponseFluentHeader interface {
	// Get retrieves the first header value associated with key.
	Get(key string) string
	// GetAll retrieves all header values associated with key.
	GetAll(key string) []string
	// Keys lists all header keys present in the response.
	Keys() []string
}

// ResponseFluentRequest exposes information about the request that generated
// this response.
type ResponseFluentRequest interface {
	// Raw returns the originating request.
	Raw() *http.Request
	// Method reports the HTTP method of the request.
	Method() string
	// URL reports the request URL.
	URL() string
	// Headers returns the headers sent with the request.
	Headers() http.Header
}

// ResponseFluentStatus reports the HTTP status code along with a rich set of
// predicates for common status checks.
type ResponseFluentStatus interface {
	// Code returns the numeric HTTP status code.
	Code() int
	// Text returns the standard HTTP status text.
	Text() string
	// Is1xxInformational reports whether the code is in the 1xx range.
	Is1xxInformational() bool
	// Is2xxSuccessful reports whether the code is in the 2xx range.
	Is2xxSuccessful() bool
	// Is3xxRedirection reports whether the code is in the 3xx range.
	Is3xxRedirection() bool
	// Is4xxClientError reports whether the code is in the 4xx range.
	Is4xxClientError() bool
	// Is5xxServerError reports whether the code is in the 5xx range.
	Is5xxServerError() bool
	// IsError reports whether the code is 4xx or 5xx.
	IsError() bool
	ResponseStatus
}

// ResponseStatus contains boolean helpers for every standard HTTP status code.
//
//nolint:interfacebloat // The interface intentionally lists many methods.
type ResponseStatus interface {
	// 1xx codes

	IsContinue() bool           // IsContinue is true for 100 Continue.
	IsSwitchingProtocols() bool // IsSwitchingProtocols is true for 101 Switching Protocols.
	IsProcessing() bool         // IsProcessing is true for 102 Processing.

	// 2xx codes

	IsOK() bool                   // IsOK is true for 200 OK.
	IsCreated() bool              // IsCreated is true for 201 Created.
	IsAccepted() bool             // IsAccepted is true for 202 Accepted.
	IsNonAuthoritatizeInfo() bool // IsNonAuthoritatizeInfo is true for 203 Non-Authoritative Information.
	IsNoContent() bool            // IsNoContent is true for 204 No Content.
	IsResetContent() bool         // IsResetContent is true for 205 Reset Content.
	IsPartialContent() bool       // IsPartialContent is true for 206 Partial Content.
	IsMultiStatus() bool          // IsMultiStatus is true for 207 Multi-Status.
	IsAlreadyReported() bool      // IsAlreadyReported is true for 208 Already Reported.
	IsIMUsed() bool               // IsIMUsed is true for 226 IM Used.

	// 3xx codes

	IsMultipleChoices() bool   // IsMultipleChoices is true for 300 Multiple Choices.
	IsMovedPermanently() bool  // IsMovedPermanently is true for 301 Moved Permanently.
	IsFound() bool             // IsFound is true for 302 Found.
	IsSeeOther() bool          // IsSeeOther is true for 303 See Other.
	IsNotModified() bool       // IsNotModified is true for 304 Not Modified.
	IsUseProxy() bool          // IsUseProxy is true for 305 Use Proxy.
	IsUnused() bool            // IsUnused is true for 306 (Unused).
	IsTemporaryRedirect() bool // IsTemporaryRedirect is true for 307 Temporary Redirect.
	IsPermanentRedirect() bool // IsPermanentRedirect is true for 308 Permanent Redirect.

	// 4xx codes

	IsBadRequest() bool                   // IsBadRequest is true for 400 Bad Request.
	IsUnauthorized() bool                 // IsUnauthorized is true for 401 Unauthorized.
	IsPaymentRequired() bool              // IsPaymentRequired is true for 402 Payment Required.
	IsForbidden() bool                    // IsForbidden is true for 403 Forbidden.
	IsNotFound() bool                     // IsNotFound is true for 404 Not Found.
	IsMethodNotAllowed() bool             // IsMethodNotAllowed is true for 405 Method Not Allowed.
	IsNotAcceptable() bool                // IsNotAcceptable is true for 406 Not Acceptable.
	IsProxyAuthRequired() bool            // IsProxyAuthRequired is true for 407 Proxy Authentication Required.
	IsRequestTimeout() bool               // IsRequestTimeout is true for 408 Request Timeout.
	IsConflict() bool                     // IsConflict is true for 409 Conflict.
	IsGone() bool                         // IsGone is true for 410 Gone.
	IsLengthRequired() bool               // IsLengthRequired is true for 411 Length Required.
	IsPreconditionFailed() bool           // IsPreconditionFailed is true for 412 Precondition Failed.
	IsRequestEntityTooLarge() bool        // IsRequestEntityTooLarge is true for 413 Content Too Large.
	IsRequestURITooLong() bool            // IsRequestURITooLong is true for 414 URI Too Long.
	IsUnsupportedMediaType() bool         // IsUnsupportedMediaType is true for 415 Unsupported Media Type.
	IsRequestedRangeNotSatisfiable() bool // IsRequestedRangeNotSatisfiable is true for 416 Range Not Satisfiable.
	IsExpectationFailed() bool            // IsExpectationFailed is true for 417 Expectation Failed.
	IsTeapot() bool                       // IsTeapot is true for 418 I'm a teapot.
	IsMisdirectedRequest() bool           // IsMisdirectedRequest is true for 421 Misdirected Request.
	IsUnprocessableEntity() bool          // IsUnprocessableEntity is true for 422 Unprocessable Entity.
	IsLocked() bool                       // IsLocked is true for 423 Locked.
	IsFailedDependency() bool             // IsFailedDependency is true for 424 Failed Dependency.
	IsTooEarly() bool                     // IsTooEarly is true for 425 Too Early.
	IsUpgradeRequired() bool              // IsUpgradeRequired is true for 426 Upgrade Required.
	IsPreconditionRequired() bool         // IsPreconditionRequired is true for 428 Precondition Required.
	IsTooManyRequests() bool              // IsTooManyRequests is true for 429 Too Many Requests.
	IsRequestHeaderFieldsTooLarge() bool  // IsRequestHeaderFieldsTooLarge is true for 431 Request Header Fields Too Large.
	IsUnavailableForLegalReasons() bool   // IsUnavailableForLegalReasons is true for 451 Unavailable For Legal Reasons.

	// 5xx codes

	IsInternalServerError() bool           // IsInternalServerError is true for 500 Internal Server Error.
	IsNotImplemented() bool                // IsNotImplemented is true for 501 Not Implemented.
	IsBadGateway() bool                    // IsBadGateway is true for 502 Bad Gateway.
	IsServiceUnavailable() bool            // IsServiceUnavailable is true for 503 Service Unavailable.
	IsGatewayTimeout() bool                // IsGatewayTimeout is true for 504 Gateway Timeout.
	IsHTTPVersionNotSupported() bool       // IsHTTPVersionNotSupported is true for 505 HTTP Version Not Supported.
	IsVariantAlsoNegotiates() bool         // IsVariantAlsoNegotiates is true for 506 Variant Also Negotiates.
	IsInsufficientStorage() bool           // IsInsufficientStorage is true for 507 Insufficient Storage.
	IsLoopDetected() bool                  // IsLoopDetected is true for 508 Loop Detected.
	IsNotExtended() bool                   // IsNotExtended is true for 510 Not Extended.
	IsNetworkAuthenticationRequired() bool // IsNetworkAuthenticationRequired is true for 511 Network Authentication Required.
}
