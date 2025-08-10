package contracts

import (
	"net/http"

	"github.com/jeanmolossi/MaiGo/pkg/maigo/header"
)

// Header is the interface that wraps the basic methods for managing HTTP headers.
//
// This wrapper provides an abstraction layer over the standard http.Header type,
// allowing for type-safe header manipulation and potential future enhancements
// without changing the public API.
//
// It enables the package to implement custom header handling logic, such as
// case-insensitive header matching or header-specific validations, while
// maintaining a consistent interface.
//
// Example:
//
//	type CustomHeader struct {
//	    header http.Header
//	}
//
//	func (c *CustomHeader) Set(key header.Type, value string) {
//	    c.header.Set(string(key), value)
//	    // Add your logic here...
//	}
type Header interface {
	// Unwrap returns the underlying header map.
	Unwrap() *http.Header
	// Get retrieves the first value associated with the key.
	Get(key header.Type) string
	// Add appends a value for the key.
	Add(key header.Type, value string)
	// Set sets the header value, replacing any existing values.
	Set(key header.Type, value string)
}
