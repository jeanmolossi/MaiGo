package maigo

import (
	"net/http"
	"sync"

	"github.com/jeanmolossi/maigo/pkg/maigo/contracts"
	"github.com/jeanmolossi/maigo/pkg/maigo/header"
	"golang.org/x/net/http/httpguts"
)

var _ contracts.Header = (*Header)(nil)

// Header wraps an http.Header map providing concurrency-safe access
// and validation of header names and values according to RFC 7230.
//
// A nil Header is treated as an empty map; methods are no-ops when the
// receiver is nil.
type Header struct {
	mu     sync.RWMutex
	header http.Header
}

// Add appends value to the field named by key. It creates the map on
// first use and silently discards invalid names or values.
func (h *Header) Add(key header.Type, value string) {
	if h == nil {
		return
	}

	if !httpguts.ValidHeaderFieldName(key.String()) || !httpguts.ValidHeaderFieldValue(value) {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.header == nil {
		h.header = make(http.Header)
	}

	h.header.Add(key.String(), value)
}

// Get retrieves the first value associated with key. It returns an
// empty string if the receiver is nil, the map is uninitialized or the
// key is invalid or absent.
func (h *Header) Get(key header.Type) string {
	if h == nil || h.header == nil {
		return ""
	}

	if !httpguts.ValidHeaderFieldName(key.String()) {
		return ""
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.header.Get(key.String())
}

// Set replaces the current value of key with value. It initializes the
// map on first use and silently ignores invalid names or values.
func (h *Header) Set(key header.Type, value string) {
	if h == nil {
		return
	}

	if !httpguts.ValidHeaderFieldName(key.String()) || !httpguts.ValidHeaderFieldValue(value) {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.header == nil {
		h.header = make(http.Header)
	}

	h.header.Set(key.String(), value)
}

// Unwrap returns a copy of the underlying header map. The caller may
// mutate the returned map without affecting the original, enabling safe
// inspection outside of the mutex.
func (h *Header) Unwrap() *http.Header {
	if h == nil {
		hdr := make(http.Header)
		return &hdr
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.header == nil {
		hdr := make(http.Header)
		return &hdr
	}

	copyHdr := make(http.Header, len(h.header))

	for k, v := range h.header {
		vv := make([]string, len(v))
		copy(vv, v)
		copyHdr[k] = vv
	}

	return &copyHdr
}

// newDefaultHTTPHeader initializes a new Header with an empty map.
func newDefaultHTTPHeader() *Header {
	return &Header{
		header: make(http.Header),
	}
}
