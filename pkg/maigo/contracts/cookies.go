package contracts

import "net/http"

// Cookies wraps basic methods for managing HTTP cookies.
//
// This wrapper provides a unified interface for cookie management, abstracting
// away the details of cookie storage and retrieval.
//
// It allows the package to implement different cookie storage strategies
// (in-memory, persistent storage) without affecting the public API. It also
// facilitates easier testing and mocking of cookie-related functionality.
//
// Example:
//
//	type CustomCookie struct {
//	    storage CookieStorage
//	}
//
//	func (c *CustomCookie) Add(cookie *http.Cookie) {
//	    c.storage.Save(cookie)
//	    // Additional logic
//	}
type Cookies interface {
	// Unwrap returns all stored cookies.
	Unwrap() []*http.Cookie
	// Get retrieves a cookie by index.
	Get(index int) *http.Cookie
	// Len reports how many cookies are stored.
	Len() int
	// Count reports how many cookies are stored.
	// Deprecated: use Len.
	Count() int
	// Add inserts a new cookie into the store.
	Add(cookie *http.Cookie)
}
