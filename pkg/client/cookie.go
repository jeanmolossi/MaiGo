package client

import (
	"net/http"

	"github.com/jeanmolossi/MaiGo/pkg/client/contracts"
)

var _ contracts.Cookies = (*Cookies)(nil)

type Cookies struct {
	cookies []*http.Cookie
}

// Add implements contracts.Cookies.
func (c *Cookies) Add(cookie *http.Cookie) {
	c.cookies = append(c.cookies, cookie)
}

// Count implements contracts.Cookies.
func (c *Cookies) Count() int {
	return len(c.cookies)
}

// Get implements contracts.Cookies.
func (c *Cookies) Get(index int) *http.Cookie {
	return c.cookies[index]
}

// Unwrap implements contracts.Cookies.
func (c *Cookies) Unwrap() []*http.Cookie {
	return c.cookies
}

// newDefaultHttpCookies initializes a new Cookies pointer.
func newDefaultHttpCookies() *Cookies {
	return &Cookies{
		cookies: make([]*http.Cookie, 0, 5), // pre-alloc 5 cookies
	}
}
