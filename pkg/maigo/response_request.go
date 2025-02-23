package maigo

import (
	"net/http"

	"github.com/jeanmolossi/MaiGo/pkg/maigo/contracts"
)

var _ contracts.ResponseFluentRequest = (*ResponseRequest)(nil)

type ResponseRequest struct {
	request *http.Request
}

// Headers implements contracts.ResponseFluentRequest.
func (r *ResponseRequest) Headers() http.Header {
	return r.request.Header
}

// Method implements contracts.ResponseFluentRequest.
func (r *ResponseRequest) Method() string {
	return r.request.Method
}

// Raw implements contracts.ResponseFluentRequest.
func (r *ResponseRequest) Raw() *http.Request {
	return r.request
}

// URL implements contracts.ResponseFluentRequest.
func (r *ResponseRequest) URL() string {
	return r.request.URL.String()
}
