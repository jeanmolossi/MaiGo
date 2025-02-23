package client

import (
	"net/url"
	"sync/atomic"

	"github.com/jeanmolossi/MaiGo/pkg/client/contracts"
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

// BaseURL for BalancedBaseURL return the next base URL in the list.
func (b *BalancedBaseURL) BaseURL() *url.URL {
	currentIndex := atomic.LoadUint32(&b.currentBaseURL)
	atomic.AddUint32(&b.currentBaseURL, 1)
	b.currentBaseURL = b.currentBaseURL % uint32(len(b.baseURLs))

	return b.baseURLs[currentIndex]
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
