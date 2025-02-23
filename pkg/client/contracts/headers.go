package contracts

import (
	"net/http"

	"github.com/jeanmolossi/MaiGo/pkg/client/header"
)

// Header is the interface that wraps the basic methods for managing HTTP headers.
//
// This wrapper provides an abstraction layer over the standard http.Header type,
// allowing for type-safe header manipulation and potential future enhancements without
// changing the public API.
//
// It enables the package to implement custom header handling logic, such as
// case-insensitive header matching or header-specific validations, while maintaining
// a consistent interface for both internal use and potential extension points.
//
// Example:
//
//	type CustomHeader struct {
//	    header http.Header
//	}
//
//	func (c *CustomHeader) Set(key header.Type, value string) {
//	    w.header.Set(string(key), value)
//	    // Add your logic here...
//	}
type Header interface {
	Unwrap() *http.Header
	Get(key header.Type) string
	Add(key header.Type, value string)
	Set(key header.Type, value string)
}
