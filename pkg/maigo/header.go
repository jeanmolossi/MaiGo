package maigo

import (
	"net/http"
	"sync"

	"github.com/jeanmolossi/maigo/pkg/maigo/contracts"
	"github.com/jeanmolossi/maigo/pkg/maigo/header"
	"golang.org/x/net/http/httpguts"
)

var _ contracts.Header = (*Header)(nil)

// Header wraps an http.Header map providing concurrency-safe access and
// validation of header names and values according to RFC 9110.
// A nil Header behaves like an empty map; all methods are no-ops on a nil
// receiver.
type Header struct {
	mu  sync.RWMutex
	hdr http.Header
}

// Add appends value to the field named by key. It creates the map on
// first use and silently discards invalid names or values.
func (h *Header) Add(key header.Type, value string) {
	if h == nil {
		return
	}

	ks := key.String()

	if !httpguts.ValidHeaderFieldName(ks) || !httpguts.ValidHeaderFieldValue(value) {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.hdr == nil {
		h.hdr = make(http.Header)
	}

	h.hdr.Add(ks, value)
}

// Get retrieves the first value associated with key. It returns an
// empty string if the receiver is nil, the map is uninitialized or the
// key is invalid or absent.
func (h *Header) Get(key header.Type) string {
	ks := key.String()

	if h == nil || h.hdr == nil {
		return ""
	}

	if !httpguts.ValidHeaderFieldName(ks) {
		return ""
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.hdr.Get(ks)
}

// Set replaces the current value of key with value. It initializes the
// map on first use and silently ignores invalid names or values.
func (h *Header) Set(key header.Type, value string) {
	if h == nil {
		return
	}

	ks := key.String()

	if !httpguts.ValidHeaderFieldName(ks) || !httpguts.ValidHeaderFieldValue(value) {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.hdr == nil {
		h.hdr = make(http.Header)
	}

	h.hdr.Set(ks, value)
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

	if h.hdr == nil {
		hdr := make(http.Header)
		return &hdr
	}

	cloned := h.hdr.Clone()

	return &cloned
}

// newDefaultHTTPHeader initializes a new Header with an empty map.
func newDefaultHTTPHeader() *Header {
	return &Header{
		hdr: make(http.Header),
	}
}
