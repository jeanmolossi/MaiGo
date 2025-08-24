package maigo

import (
	"net/http"

	"github.com/jeanmolossi/maigo/pkg/maigo/contracts"
)

var _ contracts.Cookies = (*Cookies)(nil)

// Cookies stores an in-memory collection of HTTP cookies.
//
// The zero value is ready to use. The type is not safe for concurrent use;
// callers must ensure their own synchronization if the same instance is
// accessed from multiple goroutines. The helper newDefaultHTTPCookies
// pre-allocates space for a small number of cookies.
type Cookies struct {
	cookies []*http.Cookie
}

// Add appends cookie to the collection.
//
// Nil cookies are ignored. The caller is responsible for providing a fully
// initialized *http.Cookie and for any duplicate management or validation.
func (c *Cookies) Add(cookie *http.Cookie) {
	if cookie == nil {
		return
	}

	c.cookies = append(c.cookies, cookie)
}

// Count returns the number of stored cookies.
//
// It reports the raw number of cookies and does not account for duplicates.
// Callers concerned with uniqueness must handle it themselves.
func (c *Cookies) Count() int {
	return len(c.cookies)
}

// Get retrieves the cookie at index.
//
// It returns nil if index is out of bounds. Callers must check for a nil
// result before dereferencing and ensure the index provided is intended.
func (c *Cookies) Get(index int) *http.Cookie {
	if index < 0 || index >= len(c.cookies) {
		return nil
	}

	return c.cookies[index]
}

// Unwrap returns a copy of the stored cookies.
//
// The returned slice is a new allocation, so modifying it will not affect the
// internal state. It returns nil when no cookies are stored. The caller may
// freely modify the returned slice.
func (c *Cookies) Unwrap() []*http.Cookie {
	if len(c.cookies) == 0 {
		return nil
	}

	out := make([]*http.Cookie, len(c.cookies))
	copy(out, c.cookies)

	return out
}

// newDefaultHTTPCookies creates a Cookies instance with space for a few
// cookies. The zero value of Cookies is also ready to use, so callers should
// decide whether this pre-allocation is necessary.
func newDefaultHTTPCookies() *Cookies {
	return &Cookies{
		cookies: make([]*http.Cookie, 0, 5), // pre-alloc 5 cookies
	}
}
