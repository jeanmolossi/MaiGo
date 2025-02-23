package contracts

import "net/http"

// Cookies is the interface that wraps the basic methods for managing HTTP cookies.
//
// This wrapper provides a unified interface for cookie management, abstracting
// away the details of cookie storage and retrieval.
//
// It allows the package to implement different cookie storage strategies (in-memory, persistent storage)
// without affecting the public API. It also facilitates easier testing and mocking of cookie-related funcionality.
//
// Example:
//
//	type CustomCookie struct {
//	    storage CookieStorage
//	}
//
//	func (c *CustomCookie) Add(cookie *http.Cookies) {
//	    w.storage.Save(cookie)
//	    // Additional logic
//	}
type Cookies interface {
	Unwrap() []*http.Cookie
	Get(index int) *http.Cookie
	Count() int
	Add(cookie *http.Cookie)
}
