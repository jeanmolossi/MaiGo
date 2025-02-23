package contracts

import (
	"io"
	"net/http"
)

type Response interface {
	Raw() *http.Response
	Body() ResponseFluentBody
	Header() ResponseFluentHeader
	Cookie() ResponseFluentCookie
	Request() ResponseFluentRequest
	Status() ResponseFluentStatus
}

type ResponseFluentBody interface {
	Raw() io.Closer
	Close()
	AsBytes() ([]byte, error)
	AsString() (string, error)
	AsJSON(v any) error
	AsXML(v any) error
}

type ResponseFluentCookie interface {
	GetAll() []*http.Cookie
}

type ResponseFluentHeader interface {
	Get(key string) string
	GetAll(key string) []string
	Keys() []string
}

type ResponseFluentRequest interface {
	Raw() *http.Request
	Method() string
	URL() string
	Headers() http.Header
}

type ResponseFluentStatus interface {
	Code() int
	Text() string
	Is1xxInformational() bool
	Is2xxSuccessful() bool
	Is3xxRedirection() bool
	Is4xxClientError() bool
	Is500ServerError() bool
	IsError() bool
	ResponseStatus
}

//nolint:interfacebloat // That interface should contain all codes boolean implementations
type ResponseStatus interface {
	// 1xx codes

	IsContinue() bool
	IsSwitchingProtocols() bool
	IsProcessing() bool

	// 2xx codes

	IsOK() bool
	IsCreated() bool
	IsAccepted() bool
	IsNonAuthoritatizeInfo() bool
	IsNoContent() bool
	IsResetContent() bool
	IsPartialContent() bool
	IsMultiStatus() bool
	IsAlreadyReported() bool
	IsIMUsed() bool

	// 3xx codes

	IsMultipleChoices() bool
	IsMovedPermanently() bool
	IsFound() bool
	IsSeeOther() bool
	IsNotModified() bool
	IsUseProxy() bool
	IsUnused() bool // unused
	IsTemporaryRedirect() bool
	IsPermanentRedirect() bool

	// 4xx codes

	IsBadRequest() bool
	IsUnauthorized() bool
	IsPaymentRequired() bool
	IsForbidden() bool
	IsNotFound() bool
	IsMethodNotAllowed() bool
	IsNotAcceptable() bool
	IsProxyAuthRequired() bool
	IsRequestTimeout() bool
	IsConflict() bool
	IsGone() bool
	IsLengthRequired() bool
	IsPreconditionFailed() bool
	IsRequestEntityTooLarge() bool
	IsRequestURITooLong() bool
	IsUnsupportedMediaType() bool
	IsRequestedRangeNotSatisfiable() bool
	IsExpectationFailed() bool
	IsTeapot() bool
	IsMisdirectedRequest() bool
	IsUnprocessableEntity() bool
	IsLocked() bool
	IsFailedDependency() bool
	IsTooEarly() bool
	IsUpgradeRequired() bool
	IsPreconditionRequired() bool
	IsTooManyRequests() bool
	IsRequestHeaderFieldsTooLarge() bool
	IsUnavailableForLegalReasons() bool

	// 5xx codes

	IsInternalServerError() bool
	IsNotImplemented() bool
	IsBadGateway() bool
	IsServiceUnavailable() bool
	IsGatewayTimeout() bool
	IsHTTPVersionNotSupported() bool
	IsVariantAlsoNegotiates() bool
	IsInsufficientStorage() bool
	IsLoopDetected() bool
	IsNotExtended() bool
	IsNetworkAuthenticationRequired() bool
}
