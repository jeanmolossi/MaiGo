package maigo

import (
	"net/http"

	"github.com/jeanmolossi/maigo/pkg/maigo/contracts"
)

var _ contracts.Cookies = (*Cookies)(nil)

// Cookies holds HTTP cookies in memory. The zero value is ready to use, but the
// type is not safe for concurrent use.
type Cookies struct {
	cookies []*http.Cookie
}

const defaultCookieCap = 5 // typical requests send fewer than five cookies

// isValidCookieName reports whether name consists solely of tchar characters as
// defined by RFC 6265 §4.1.1 and RFC 9110 §5.6.2. The string must be non-empty
// and contain only characters in !#$%&'*+-.^_|~0-9A-Za-z`.
func isValidCookieName(name string) bool {
	if name == "" {
		return false
	}

	for i := 0; i < len(name); i++ {
		c := name[i]

		switch {
		case '0' <= c && c <= '9':
			continue
		case 'A' <= c && c <= 'Z':
			continue
		case 'a' <= c && c <= 'z':
			continue
		case c == '!' || c == '#' || c == '$' || c == '%' || c == '&' ||
			c == '\'' || c == '*' || c == '+' || c == '-' || c == '.' ||
			c == '^' || c == '_' || c == '`' || c == '|' || c == '~':
			continue
		default:
			return false
		}
	}

	return true
}

// Add clones cookie and appends it. Nil receiver, nil cookie, or invalid Name
// (as determined by isValidCookieName) are ignored. Duplicates are caller
// responsibility.
func (c *Cookies) Add(cookie *http.Cookie) {
	if c == nil || cookie == nil || !isValidCookieName(cookie.Name) {
		return
	}

	c.cookies = append(c.cookies, cloneCookie(cookie))
}

// Count reports how many cookies are stored.
//
// Deprecated: use Len. Count will be removed in v2.
func (c *Cookies) Count() int {
	if c == nil {
		return 0
	}

	return len(c.cookies)
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
