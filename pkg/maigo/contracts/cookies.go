package contracts

import "net/http"

// Cookies manages HTTP cookies in memory.
//
// Implementations are not safe for concurrent use.
type Cookies interface {
	// Unwrap returns deep copies of all stored cookies.
	// It returns nil when no cookies are stored, and a nil *Cookies value also unwraps to nil.
	Unwrap() []*http.Cookie
	// Get returns a deep copy of the cookie at index.
	// It returns nil if index is out of range.
	Get(index int) *http.Cookie
	// Len reports how many cookies are stored.
	Len() int
	// Count reports how many cookies are stored.
	// Deprecated: use Len. Count will be removed in v2.
	Count() int
	// Add creates a copy of cookie, trims the Name on the copy using strings.TrimSpace,
	// and stores that copy. The input cookie is not mutated. Nil or blank-name cookies
	// (after trimming) are ignored.
	Add(cookie *http.Cookie)
}
