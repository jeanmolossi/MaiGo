package maigo

import (
	"net/http"
	"strings"

	"github.com/jeanmolossi/maigo/pkg/maigo/contracts"
)

var _ contracts.Cookies = (*Cookies)(nil)

// Cookies holds HTTP cookies in memory. The zero value is ready to use, but the
// type is not safe for concurrent use.
type Cookies struct {
	cookies []*http.Cookie
}

const defaultCookieCap = 5 // typical requests send fewer than five cookies

var tcharTable = func() [128]bool {
	var tbl [128]bool
	for c := '0'; c <= '9'; c++ {
		tbl[c] = true
	}
	for c := 'A'; c <= 'Z'; c++ {
		tbl[c] = true
	}
	for c := 'a'; c <= 'z'; c++ {
		tbl[c] = true
	}
	for _, c := range "!#$%&'*+-.^_|~`" {
		tbl[c] = true
	}

	return tbl
}()

// isValidCookieName reports whether name consists solely of tchar characters as
// defined by RFC 6265 §4.1.1 and RFC 9110 §5.6.2. The string must be non-empty
// and contain only characters in !#$%&'*+-.^_|~0-9A-Za-z`.
func isValidCookieName(name string) bool {
	if name == "" {
		return false
	}

	for i := 0; i < len(name); i++ {
		c := name[i]
		if c >= 128 || !tcharTable[c] {
			return false
		}
	}

	return true
}

// Add creates a copy of cookie, trims the Name on the copy using
// strings.TrimSpace, and appends it. Nil receiver or nil cookie are ignored.
// The Name must be non-empty after trimming and pass isValidCookieName. The
// stored cookie's Name is the trimmed value; the input cookie is not mutated.
// Handling duplicates remains the caller's responsibility.
func (c *Cookies) Add(cookie *http.Cookie) {
	if c == nil || cookie == nil {
		return
	}

	name := strings.TrimSpace(cookie.Name)
	if name == "" || !isValidCookieName(name) {
		return
	}

	clone := cloneCookie(cookie)
	clone.Name = name
	c.cookies = append(c.cookies, clone)
}

// Count reports how many cookies are stored.
//
// Deprecated: use Len. Count will be removed in v2.
func (c *Cookies) Count() int {
	return c.Len()
}

// Len reports how many cookies are stored.
func (c *Cookies) Len() int {
	if c == nil {
		return 0
	}

	return len(c.cookies)
}

// Get returns a clone of the cookie at index or nil if out of range.
func (c *Cookies) Get(index int) *http.Cookie {
	if c == nil || index < 0 || index >= len(c.cookies) {
		return nil
	}

	return cloneCookie(c.cookies[index])
}

// Unwrap returns deep copies of all stored cookies or nil when empty.
func (c *Cookies) Unwrap() []*http.Cookie {
	if c == nil || len(c.cookies) == 0 {
		return nil
	}

	out := make([]*http.Cookie, len(c.cookies))
	for i, ck := range c.cookies {
		out[i] = cloneCookie(ck)
	}

	return out
}

// newDefaultHTTPCookies creates a Cookies instance with room for a few cookies.
func newDefaultHTTPCookies() *Cookies {
	return newCookiesWithCapacity(defaultCookieCap)
}

// newCookiesWithCapacity returns a Cookies instance pre-allocated to capacity.
func newCookiesWithCapacity(capacity int) *Cookies {
	if capacity < 0 {
		capacity = 0
	}

	return &Cookies{cookies: make([]*http.Cookie, 0, capacity)}
}

func cloneCookie(src *http.Cookie) *http.Cookie {
	if src == nil {
		return nil
	}

	dst := new(http.Cookie)
	*dst = *src

	if src.Unparsed != nil {
		up := make([]string, len(src.Unparsed))
		copy(up, src.Unparsed)
		dst.Unparsed = up
	}

	return dst
}
