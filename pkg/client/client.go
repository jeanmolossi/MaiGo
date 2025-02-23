package client

import (
	"net/http"
	"time"

	"github.com/jeanmolossi/MaiGo/pkg/client/contracts"
)

// interface implementation type check.
var _ contracts.HTTPClient = (*Client)(nil)

type Client struct {
	client *http.Client
}

// Do implements contracts.HTTPClient.
func (d *Client) Do(r *http.Request) (*http.Response, error) {
	//nolint:wrapcheck // that was intentional
	return d.client.Do(r)
}

// SetFollowRedirects implements contracts.HTTPClient.
func (d *Client) SetFollowRedirects(follow bool) {
	if !follow {
		d.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
}

// SetTimeout implements contracts.HTTPClient.
func (d *Client) SetTimeout(timeout time.Duration) {
	d.client.Timeout = timeout
}

// SetTransport implements contracts.HTTPClient.
func (d *Client) SetTransport(rt http.RoundTripper) {
	d.client.Transport = rt
}

// Timeout implements contracts.HTTPClient.
func (d *Client) Timeout() time.Duration {
	return d.client.Timeout
}

// Transport implements contracts.HTTPClient.
func (d *Client) Transport() http.RoundTripper {
	return d.client.Transport
}

// newDefaultHTTPClient initialized a new DefaultHttpClient.
func newDefaultHTTPClient() *Client {
	return &Client{
		client: &http.Client{},
	}
}
