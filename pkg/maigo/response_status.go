package maigo

import (
	"fmt"
	"net/http"

	"github.com/jeanmolossi/MaiGo/pkg/maigo/contracts"
)

var _ contracts.ResponseFluentStatus = (*ResponseStatus)(nil)

type ResponseStatus struct {
	response *http.Response
}

// Code check if status is xxx.
func (r *ResponseStatus) Code() int {
	return r.response.StatusCode
}

// Text check if status is xxx.
func (r *ResponseStatus) Text() string {
	return fmt.Sprintf(
		"[%d] %s",
		r.Code(),
		http.StatusText(r.Code()),
	)
}

// Is1xxInformational check if status is xxx.
func (r *ResponseStatus) Is1xxInformational() bool {
	return r.response.StatusCode >= 100 && r.response.StatusCode < 200
}

// Is2xxSuccessful check if status is xxx.
func (r *ResponseStatus) Is2xxSuccessful() bool {
	return r.response.StatusCode >= 200 && r.response.StatusCode < 300
}

// Is3xxRedirection check if status is xxx.
func (r *ResponseStatus) Is3xxRedirection() bool {
	return r.response.StatusCode >= 300 && r.response.StatusCode < 400
}

// Is4xxClientError check if status is xxx.
func (r *ResponseStatus) Is4xxClientError() bool {
	return r.response.StatusCode >= 400 && r.response.StatusCode < 500
}

// Is5xxServerError check if status is xxx.
func (r *ResponseStatus) Is5xxServerError() bool {
	return r.response.StatusCode >= 500 && r.response.StatusCode < 600
}

// IsContinue check if status is 100.
func (r *ResponseStatus) IsContinue() bool {
	return r.response.StatusCode == http.StatusContinue
}

// IsSwitchingProtocols check if status is 101.
func (r *ResponseStatus) IsSwitchingProtocols() bool {
	return r.response.StatusCode == http.StatusSwitchingProtocols
}

// IsProcessing check if status is 102.
func (r *ResponseStatus) IsProcessing() bool {
	return r.response.StatusCode == http.StatusProcessing
}

// IsEarlyHints check if status is 103.
func (r *ResponseStatus) IsEarlyHints() bool {
	return r.response.StatusCode == http.StatusEarlyHints
}

// ----------------------------------------------------
//
// 2xx Codes
//
// ----------------------------------------------------

// IsOK check if status is 200.
func (r *ResponseStatus) IsOK() bool {
	return r.response.StatusCode == http.StatusOK
}

// IsCreated check if status is 201.
func (r *ResponseStatus) IsCreated() bool {
	return r.response.StatusCode == http.StatusCreated
}

// IsAccepted check if status is 202.
func (r *ResponseStatus) IsAccepted() bool {
	return r.response.StatusCode == http.StatusAccepted
}

// IsNonAuthoritatizeInfo check if status is 203.
func (r *ResponseStatus) IsNonAuthoritatizeInfo() bool {
	return r.response.StatusCode == http.StatusNonAuthoritativeInfo
}

// IsNoContent check if status is 204.
func (r *ResponseStatus) IsNoContent() bool {
	return r.response.StatusCode == http.StatusNoContent
}

// IsResetContent check if status is 205.
func (r *ResponseStatus) IsResetContent() bool {
	return r.response.StatusCode == http.StatusResetContent
}

// IsPartialContent check if status is 206.
func (r *ResponseStatus) IsPartialContent() bool {
	return r.response.StatusCode == http.StatusPartialContent
}

// IsMultiStatus check if status is 207.
func (r *ResponseStatus) IsMultiStatus() bool {
	return r.response.StatusCode == http.StatusMultiStatus
}

// IsAlreadyReported check if status is 208.
func (r *ResponseStatus) IsAlreadyReported() bool {
	return r.response.StatusCode == http.StatusAlreadyReported
}

// IsIMUsed check if status is 226.
func (r *ResponseStatus) IsIMUsed() bool {
	return r.response.StatusCode == http.StatusIMUsed
}

// ----------------------------------------------------
//
// 3xx Codes
//
// ----------------------------------------------------

// IsMultipleChoices check if status is 300.
func (r *ResponseStatus) IsMultipleChoices() bool {
	return r.response.StatusCode == http.StatusMultipleChoices
}

// IsMovedPermanently check if status is 301.
func (r *ResponseStatus) IsMovedPermanently() bool {
	return r.response.StatusCode == http.StatusMovedPermanently
}

// IsFound check if status is 302.
func (r *ResponseStatus) IsFound() bool {
	return r.response.StatusCode == http.StatusFound
}

// IsSeeOther check if status is 303.
func (r *ResponseStatus) IsSeeOther() bool {
	return r.response.StatusCode == http.StatusSeeOther
}

// IsNotModified check if status is 304.
func (r *ResponseStatus) IsNotModified() bool {
	return r.response.StatusCode == http.StatusNotModified
}

// IsUseProxy check if status is 305.
func (r *ResponseStatus) IsUseProxy() bool {
	return r.response.StatusCode == http.StatusUseProxy
}

// IsUnused check if status is 306.
//
// RFC 9110, 15.4.7 (Unused)
//
// The 306 status code was defined in a previous version of this specification,
// is no longer used, and the code is reserved.
//
// See: https://www.rfc-editor.org/rfc/rfc9110.html#section-15.4.7-1
func (r *ResponseStatus) IsUnused() bool {
	unusedStatus := 306
	return r.response.StatusCode == unusedStatus
}

// IsTemporaryRedirect check if status is 307.
func (r *ResponseStatus) IsTemporaryRedirect() bool {
	return r.response.StatusCode == http.StatusTemporaryRedirect
}

// IsPermanentRedirect check if status is 308.
func (r *ResponseStatus) IsPermanentRedirect() bool {
	return r.response.StatusCode == http.StatusPermanentRedirect
}

// ----------------------------------------------------
//
// 4xx Codes
//
// ----------------------------------------------------

// IsBadRequest check if status is 400.
func (r *ResponseStatus) IsBadRequest() bool {
	return r.response.StatusCode == http.StatusBadRequest
}

// IsUnauthorized check if status is 401.
func (r *ResponseStatus) IsUnauthorized() bool {
	return r.response.StatusCode == http.StatusUnauthorized
}

// IsPaymentRequired check if status is 402.
func (r *ResponseStatus) IsPaymentRequired() bool {
	return r.response.StatusCode == http.StatusPaymentRequired
}

// IsForbidden check if status is 403.
func (r *ResponseStatus) IsForbidden() bool {
	return r.response.StatusCode == http.StatusForbidden
}

// IsNotFound check if status is 404.
func (r *ResponseStatus) IsNotFound() bool {
	return r.response.StatusCode == http.StatusNotFound
}

// IsMethodNotAllowed check if status is 405.
func (r *ResponseStatus) IsMethodNotAllowed() bool {
	return r.response.StatusCode == http.StatusMethodNotAllowed
}

// IsNotAcceptable check if status is 406.
func (r *ResponseStatus) IsNotAcceptable() bool {
	return r.response.StatusCode == http.StatusNotAcceptable
}

// IsProxyAuthRequired check if status is 407.
func (r *ResponseStatus) IsProxyAuthRequired() bool {
	return r.response.StatusCode == http.StatusProxyAuthRequired
}

// IsRequestTimeout check if status is 408.
func (r *ResponseStatus) IsRequestTimeout() bool {
	return r.response.StatusCode == http.StatusRequestTimeout
}

// IsConflict check if status is 409.
func (r *ResponseStatus) IsConflict() bool {
	return r.response.StatusCode == http.StatusConflict
}

// IsGone check if status is 410.
func (r *ResponseStatus) IsGone() bool {
	return r.response.StatusCode == http.StatusGone
}

// IsLengthRequired check if status is 411.
func (r *ResponseStatus) IsLengthRequired() bool {
	return r.response.StatusCode == http.StatusLengthRequired
}

// IsPreconditionFailed check if status is 412.
func (r *ResponseStatus) IsPreconditionFailed() bool {
	return r.response.StatusCode == http.StatusPreconditionFailed
}

// IsRequestEntityTooLarge check if status is 413.
func (r *ResponseStatus) IsRequestEntityTooLarge() bool {
	return r.response.StatusCode == http.StatusRequestEntityTooLarge
}

// IsRequestURITooLong check if status is 414.
func (r *ResponseStatus) IsRequestURITooLong() bool {
	return r.response.StatusCode == http.StatusRequestURITooLong
}

// IsUnsupportedMediaType check if status is 415.
func (r *ResponseStatus) IsUnsupportedMediaType() bool {
	return r.response.StatusCode == http.StatusUnsupportedMediaType
}

// IsRequestedRangeNotSatisfiable check if status is 416.
func (r *ResponseStatus) IsRequestedRangeNotSatisfiable() bool {
	return r.response.StatusCode == http.StatusRequestedRangeNotSatisfiable
}

// IsExpectationFailed check if status is 417.
func (r *ResponseStatus) IsExpectationFailed() bool {
	return r.response.StatusCode == http.StatusExpectationFailed
}

// IsTeapot check if status is 418.
func (r *ResponseStatus) IsTeapot() bool {
	return r.response.StatusCode == http.StatusTeapot
}

// IsMisdirectedRequest check if status is 421.
func (r *ResponseStatus) IsMisdirectedRequest() bool {
	return r.response.StatusCode == http.StatusMisdirectedRequest
}

// IsUnprocessableEntity check if status is 422.
func (r *ResponseStatus) IsUnprocessableEntity() bool {
	return r.response.StatusCode == http.StatusUnprocessableEntity
}

// IsLocked check if status is 423.
func (r *ResponseStatus) IsLocked() bool {
	return r.response.StatusCode == http.StatusLocked
}

// IsFailedDependency check if status is 424.
func (r *ResponseStatus) IsFailedDependency() bool {
	return r.response.StatusCode == http.StatusFailedDependency
}

// IsTooEarly check if status is 425.
func (r *ResponseStatus) IsTooEarly() bool {
	return r.response.StatusCode == http.StatusTooEarly
}

// IsUpgradeRequired check if status is 426.
func (r *ResponseStatus) IsUpgradeRequired() bool {
	return r.response.StatusCode == http.StatusUpgradeRequired
}

// IsPreconditionRequired check if status is 428.
func (r *ResponseStatus) IsPreconditionRequired() bool {
	return r.response.StatusCode == http.StatusPreconditionRequired
}

// IsTooManyRequests check if status is 429.
func (r *ResponseStatus) IsTooManyRequests() bool {
	return r.response.StatusCode == http.StatusTooManyRequests
}

// IsRequestHeaderFieldsTooLarge check if status is 431.
func (r *ResponseStatus) IsRequestHeaderFieldsTooLarge() bool {
	return r.response.StatusCode == http.StatusRequestHeaderFieldsTooLarge
}

// IsUnavailableForLegalReasons check if status is 451.
func (r *ResponseStatus) IsUnavailableForLegalReasons() bool {
	return r.response.StatusCode == http.StatusUnavailableForLegalReasons
}

// ----------------------------------------------------
//
// 5xx Codes
//
// ----------------------------------------------------

// IsInternalServerError check if status is 500.
func (r *ResponseStatus) IsInternalServerError() bool {
	return r.response.StatusCode == http.StatusInternalServerError
}

// IsNotImplemented check if status is 501.
func (r *ResponseStatus) IsNotImplemented() bool {
	return r.response.StatusCode == http.StatusNotImplemented
}

// IsBadGateway check if status is 502.
func (r *ResponseStatus) IsBadGateway() bool {
	return r.response.StatusCode == http.StatusBadGateway
}

// IsServiceUnavailable check if status is 503.
func (r *ResponseStatus) IsServiceUnavailable() bool {
	return r.response.StatusCode == http.StatusServiceUnavailable
}

// IsGatewayTimeout check if status is 504.
func (r *ResponseStatus) IsGatewayTimeout() bool {
	return r.response.StatusCode == http.StatusGatewayTimeout
}

// IsHTTPVersionNotSupported check if status is 505.
func (r *ResponseStatus) IsHTTPVersionNotSupported() bool {
	return r.response.StatusCode == http.StatusHTTPVersionNotSupported
}

// IsVariantAlsoNegotiates check if status is 506.
func (r *ResponseStatus) IsVariantAlsoNegotiates() bool {
	return r.response.StatusCode == http.StatusVariantAlsoNegotiates
}

// IsInsufficientStorage check if status is 507.
func (r *ResponseStatus) IsInsufficientStorage() bool {
	return r.response.StatusCode == http.StatusInsufficientStorage
}

// IsLoopDetected check if status is 508.
func (r *ResponseStatus) IsLoopDetected() bool {
	return r.response.StatusCode == http.StatusLoopDetected
}

// IsNotExtended check if status is 510.
func (r *ResponseStatus) IsNotExtended() bool {
	return r.response.StatusCode == http.StatusNotExtended
}

// IsNetworkAuthenticationRequired check if status is 511.
func (r *ResponseStatus) IsNetworkAuthenticationRequired() bool {
	return r.response.StatusCode == http.StatusNetworkAuthenticationRequired
}

// IsError check if status is xxx.
func (r *ResponseStatus) IsError() bool {
	return r.Is4xxClientError() || r.Is5xxServerError()
}
