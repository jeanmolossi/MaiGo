package maigo

import (
	"net/http"
	"strings"

	"github.com/jeanmolossi/maigo/pkg/maigo/contracts"
)

var _ contracts.Cookies = (*Cookies)(nil)

// Cookies stores an in-memory collection of HTTP cookies. The zero value is
// ready to use; the type is not safe for concurrent use—callers must
// synchronize access when sharing the same instance across goroutines. The
// helper newDefaultHTTPCookies pre-allocates space for a small number of
// cookies.
type Cookies struct {
	cookies []*http.Cookie
}

const defaultCookieCap = 5 // typical requests send fewer than five cookies

// Add appends cookie to the collection.
//
// Nil cookies or cookies whose Name is empty after trimming whitespace with
// strings.TrimSpace are ignored. The Name is checked with TrimSpace but not
// mutated, so callers must supply a non-blank Name after trimming and remain
// responsible for duplicate management or any additional validation. Add
// stores the caller's pointer; mutating the cookie after adding will affect
// the stored value.
func (c *Cookies) Add(cookie *http.Cookie) {
	if cookie == nil || strings.TrimSpace(cookie.Name) == "" {
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
// The returned pointer aliases internal state; mutating it updates the
// underlying cookie stored by Cookies. Callers needing an independent copy
// should use Unwrap. It returns nil if index is out of bounds, so callers must
// check for nil before dereferencing and ensure the index is intended.
func (c *Cookies) Get(index int) *http.Cookie {
	if index < 0 || index >= len(c.cookies) {
		return nil
	}

	return c.cookies[index]
}

// Unwrap returns a deep copy of the stored cookies.
//
// The returned slice and each *http.Cookie are new allocations, so modifying
// them does not affect the internal state. It returns nil when no cookies are
// stored.
func (c *Cookies) Unwrap() []*http.Cookie {
	if len(c.cookies) == 0 {
		return nil
	}

	out := make([]*http.Cookie, len(c.cookies))

	for i, ck := range c.cookies {
		if ck == nil { // defensive; Add ignores nil cookies
			continue
		}

		clone := new(http.Cookie)
		*clone = *ck

		if len(ck.Unparsed) > 0 {
			up := make([]string, len(ck.Unparsed))
			copy(up, ck.Unparsed)
			clone.Unparsed = up
		}

		out[i] = clone
	}

	return out
}

// newDefaultHTTPCookies creates a Cookies instance with space for a few
// cookies. The zero value of Cookies is also ready to use, so callers should
// decide whether this pre-allocation is necessary.
func newDefaultHTTPCookies() *Cookies {
	return &Cookies{
		cookies: make([]*http.Cookie, 0, defaultCookieCap), // pre-alloc
	}
}
