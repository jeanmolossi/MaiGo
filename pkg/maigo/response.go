package maigo

import (
	"net/http"

	"github.com/jeanmolossi/maigo/pkg/maigo/contracts"
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
	return r.body
}

// Cookie implements contracts.Response.
func (r *Response) Cookie() contracts.ResponseFluentCookie {
	return r.cookie
}

// Header implements contracts.Response.
func (r *Response) Header() contracts.ResponseFluentHeader {
	return r.header
}

// Request implements contracts.Response.
func (r *Response) Request() contracts.ResponseFluentRequest {
	return r.request
}

// Status implements contracts.Response.
func (r *Response) Status() contracts.ResponseFluentStatus {
	return r.status
}

func (r *Response) Raw() *http.Response {
	return r.raw
}

func newResponse(response *http.Response) *Response {
	return &Response{
		raw: response,
		// Fluent API
		body: &ResponseBody{
			body: newUnbufferedBody(response.Body),
		},
		cookie: &ResponseCookie{
			cookies: response.Cookies(),
		},
		header: &ResponseHeader{
			header: response.Header,
		},
		request: &ResponseRequest{
			request: response.Request,
		},
		status: &ResponseStatus{
			response: response,
		},
	}
}
