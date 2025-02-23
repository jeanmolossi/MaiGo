package maigo

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/jeanmolossi/MaiGo/pkg/maigo/contracts"
)

var _ contracts.BuilderHTTPClientConfig[contracts.ClientBuilder] = (*ClientConfigBuilder)(nil)

type ClientConfigBuilder struct {
	parent *ClientBuilder
}

func (b *ClientBuilder) Config() *ClientConfigBuilder {
	return &ClientConfigBuilder{parent: b}
}

// SetCustomHTTPClient implements contracts.BuilderHTTPClientConfig.
func (c *ClientConfigBuilder) SetCustomHTTPClient(httpClient contracts.HTTPClient) contracts.ClientBuilder {
	c.parent.client.SetHttpClient(httpClient)
	return c.parent
}

// SetCustomTransport implements contracts.BuilderHTTPClientConfig.
func (c *ClientConfigBuilder) SetCustomTransport(transport http.RoundTripper) contracts.ClientBuilder {
	c.parent.client.HttpClient().SetTransport(transport)
	return c.parent
}

// SetFollowRedirects implements contracts.BuilderHTTPClientConfig.
func (c *ClientConfigBuilder) SetFollowRedirects(follow bool) contracts.ClientBuilder {
	c.parent.client.HttpClient().SetFollowRedirects(follow)
	return c.parent
}

// SetProxy implements contracts.BuilderHTTPClientConfig.
func (c *ClientConfigBuilder) SetProxy(proxyURL string) contracts.ClientBuilder {
	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		c.parent.client.Validations().Add(
			errors.Join(ErrParseProxyURL, err),
		)

		return c.parent
	}

	if transport, ok := c.parent.client.HttpClient().Transport().(*http.Transport); ok {
		transport.Proxy = http.ProxyURL(parsedURL)
	} else {
		c.parent.client.HttpClient().SetTransport(&http.Transport{
			Proxy: http.ProxyURL(parsedURL),
		})
	}

	return c.parent
}

// SetTimeout implements contracts.BuilderHTTPClientConfig.
func (c *ClientConfigBuilder) SetTimeout(duration time.Duration) contracts.ClientBuilder {
	c.parent.client.HttpClient().SetTimeout(duration)
	return c.parent
}
