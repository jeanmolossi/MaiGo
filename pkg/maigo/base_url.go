package maigo

import (
	"net/url"
	"sync/atomic"

	"github.com/jeanmolossi/maigo/pkg/maigo/contracts"
)

type (
	// DefaultBaseURL implements contracts.ConfigBaseURL interface and provides a single base URL.
	DefaultBaseURL struct {
		baseURL *url.URL
	}

	// BalancedBaseURL implements contracts.ConfigBaseURL interface and provides a load balancing.
	BalancedBaseURL struct {
		baseURLs       []*url.URL
		currentBaseURL uint32
	}
)

// test implementations.
var (
	_ contracts.ConfigBaseURL = (*DefaultBaseURL)(nil)
	_ contracts.ConfigBaseURL = (*BalancedBaseURL)(nil)
)

// BaseURL for DefaultBaseURL return the base URL.
func (d *DefaultBaseURL) BaseURL() *url.URL {
	return d.baseURL
}

// BaseURL for BalancedBaseURL returns the next base URL in the list.
// It is safe for concurrent use and for zero or single URLs.
func (b *BalancedBaseURL) BaseURL() *url.URL {
	l := len(b.baseURLs)
	switch l {
	case 0:
		return nil
	case 1:
		return b.baseURLs[0]
	}

	idx := atomic.AddUint32(&b.currentBaseURL, 1) - 1

	return b.baseURLs[idx%uint32(l)]
}

// newDefaultBaseURL initializes a new DefaultBaseURL with a given base URL.
func newDefaultBaseURL(baseURL *url.URL) *DefaultBaseURL {
	return &DefaultBaseURL{
		baseURL: baseURL,
	}
}

// newBalancedBaseURL initializes a new BalancedBaseURL with a given base URLs.
func newBalancedBaseURL(baseURLs []*url.URL) *BalancedBaseURL {
	return &BalancedBaseURL{
		baseURLs: baseURLs,
	}
}
