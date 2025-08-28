package contracts

import (
	"net/http"

	"github.com/jeanmolossi/maigo/pkg/maigo/header"
)

// Header defines a minimal, concurrency-safe API for manipulating HTTP
// headers. Implementations should validate field names and values according
// to RFC 7230 and may silently discard invalid input.
type Header interface {
	// Unwrap returns a copy of the header map for read-only inspection.
	Unwrap() *http.Header
	// Get retrieves the first value associated with key, or an empty string
	// if the key is absent.
	Get(key header.Type) string
	// Add appends value to key, preserving existing values.
	Add(key header.Type, value string)
	// Set replaces any existing values of key with value.
	Set(key header.Type, value string)
}
