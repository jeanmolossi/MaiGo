package client

import (
	"net/http"

	"github.com/jeanmolossi/MaiGo/pkg/client/contracts"
)

var _ contracts.ResponseFluentHeader = (*ResponseHeader)(nil)

type ResponseHeader struct {
	header http.Header
}

// Get implements contracts.ResponseFluentHeader.
func (r *ResponseHeader) Get(key string) string {
	return r.header.Get(key)
}

// GetAll implements contracts.ResponseFluentHeader.
func (r *ResponseHeader) GetAll(key string) []string {
	return r.header[key]
}

// Keys implements contracts.ResponseFluentHeader.
func (r *ResponseHeader) Keys() []string {
	keys := make([]string, 0, len(r.header))
	for k := range r.header {
		keys = append(keys, k)
	}

	return keys
}
