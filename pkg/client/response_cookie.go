package client

import (
	"net/http"

	"github.com/jeanmolossi/MaiGo/pkg/client/contracts"
)

var _ contracts.ResponseFluentCookie = (*ResponseCookie)(nil)

type ResponseCookie struct {
	cookies []*http.Cookie
}

// GetAll implements contracts.ResponseFluentCookie.
func (r *ResponseCookie) GetAll() []*http.Cookie {
	return r.cookies
}
