package client

import (
	"net/http"

	"github.com/jeanmolossi/MaiGo/pkg/client/contracts"
)

var _ contracts.Response = (*Response)(nil)

type Response struct {
	raw *http.Response

	// Fluent API
	body    contracts.ResponseFluentBody
	cookie  contracts.ResponseFluentCookie
	header  contracts.ResponseFluentHeader
	request contracts.ResponseFluentRequest
	status  contracts.ResponseFluentStatus
}

// Body implements contracts.Response.
func (r *Response) Body() contracts.ResponseFluentBody {
	panic("unimplemented")
}

// Cookie implements contracts.Response.
func (r *Response) Cookie() contracts.ResponseFluentCookie {
	panic("unimplemented")
}

// Header implements contracts.Response.
func (r *Response) Header() contracts.ResponseFluentHeader {
	panic("unimplemented")
}

// Request implements contracts.Response.
func (r *Response) Request() contracts.ResponseFluentRequest {
	panic("unimplemented")
}

// Status implements contracts.Response.
func (r *Response) Status() contracts.ResponseFluentStatus {
	panic("unimplemented")
}

func (r *Response) Raw() *http.Response {
	return r.raw
}

func newResponse(response *http.Response) *Response {
	return &Response{
		raw: response,
		// Fluent API
		body:    nil,
		cookie:  nil,
		header:  nil,
		request: nil,
		status:  nil,
	}
}
