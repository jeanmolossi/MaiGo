package maigo

import (
	"net/http"

	"github.com/jeanmolossi/maigo/pkg/maigo/contracts"
)

var _ contracts.ResponseFluentCookie = (*ResponseCookie)(nil)

type ResponseCookie struct {
	cookies []*http.Cookie
}

// GetAll implements contracts.ResponseFluentCookie.
func (r *ResponseCookie) GetAll() []*http.Cookie {
	return r.cookies
}
