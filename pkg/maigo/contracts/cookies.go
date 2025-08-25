package contracts

import "net/http"

// Cookies manages HTTP cookies in memory.
//
// Implementations are not safe for concurrent use.
type Cookies interface {
	// Unwrap returns deep copies of all stored cookies.
	// It returns nil if the collection is empty; a zero-value Cookies unwraps to nil.
	Unwrap() []*http.Cookie
	// Get returns a copy of the cookie at index or nil if out of range.
	Get(index int) *http.Cookie
	// Len reports how many cookies are stored.
	Len() int
	// Count reports how many cookies are stored.
	// Deprecated: use Len. Count will be removed in v2.
	Count() int
	// Add stores a copy of cookie. Nil or blank-name cookies are ignored.
	Add(cookie *http.Cookie)
}
