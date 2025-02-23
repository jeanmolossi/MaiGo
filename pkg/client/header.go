package client

import (
	"net/http"

	"github.com/jeanmolossi/MaiGo/pkg/client/contracts"
	"github.com/jeanmolossi/MaiGo/pkg/client/header"
)

var _ contracts.Header = (*Header)(nil)

type Header struct {
	header *http.Header
}

// Add implements contracts.Header.
func (h *Header) Add(key header.Type, value string) {
	h.header.Add(key.String(), value)
}

// Get implements contracts.Header.
func (h *Header) Get(key header.Type) string {
	return h.header.Get(key.String())
}

// Set implements contracts.Header.
func (h *Header) Set(key header.Type, value string) {
	h.header.Set(key.String(), value)
}

// Unwrap implements contracts.Header.
func (h *Header) Unwrap() *http.Header {
	return h.header
}

// newDefaultHTTPHeader initialized a new HTTPHeader.
func newDefaultHTTPHeader() *Header {
	return &Header{
		header: &http.Header{},
	}
}
