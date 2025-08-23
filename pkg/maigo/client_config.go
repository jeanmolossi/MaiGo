package maigo

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/jeanmolossi/maigo/pkg/maigo/contracts"
	"github.com/jeanmolossi/maigo/pkg/maigo/method"
)

var (
	_ contracts.ClientConfig      = (*ClientConfigBase)(nil)
	_ contracts.ClientHTTPMethods = (*ClientConfigBase)(nil)
)

// ClientConfigBase serves as the main entrypoint to configure HTTP client.
type ClientConfigBase struct {
	httpClient  contracts.HTTPClient
	httpHeader  contracts.Header
	httpCookie  contracts.Cookies
	validations contracts.Validations

	contracts.ConfigBaseURL
}

// CONNECT implements contracts.ClientHTTPMethods.
func (c *ClientConfigBase) CONNECT(path string) contracts.RequestBuilder {
	return newRequest(c, method.CONNECT, path)
}

// DELETE implements contracts.ClientHTTPMethods.
func (c *ClientConfigBase) DELETE(path string) contracts.RequestBuilder {
	return newRequest(c, method.DELETE, path)
}

// GET implements contracts.ClientHTTPMethods.
func (c *ClientConfigBase) GET(path string) contracts.RequestBuilder {
	return newRequest(c, method.GET, path)
}

// HEAD implements contracts.ClientHTTPMethods.
func (c *ClientConfigBase) HEAD(path string) contracts.RequestBuilder {
	return newRequest(c, method.HEAD, path)
}

// OPTIONS implements contracts.ClientHTTPMethods.
func (c *ClientConfigBase) OPTIONS(path string) contracts.RequestBuilder {
	return newRequest(c, method.OPTIONS, path)
}

// PATCH implements contracts.ClientHTTPMethods.
func (c *ClientConfigBase) PATCH(path string) contracts.RequestBuilder {
	return newRequest(c, method.PATCH, path)
}

// POST implements contracts.ClientHTTPMethods.
func (c *ClientConfigBase) POST(path string) contracts.RequestBuilder {
	return newRequest(c, method.POST, path)
}

// PUT implements contracts.ClientHTTPMethods.
func (c *ClientConfigBase) PUT(path string) contracts.RequestBuilder {
	return newRequest(c, method.PUT, path)
}

// TRACE implements contracts.ClientHTTPMethods.
func (c *ClientConfigBase) TRACE(path string) contracts.RequestBuilder {
	return newRequest(c, method.TRACE, path)
}

// Cookies implements contracts.ClientConfig.
func (c *ClientConfigBase) Cookies() contracts.Cookies {
	return c.httpCookie
}

// Header implements contracts.ClientConfig.
func (c *ClientConfigBase) Header() contracts.Header {
	return c.httpHeader
}

// HttpClient implements contracts.ClientConfig.
func (c *ClientConfigBase) HttpClient() contracts.HTTPClient {
	return c.httpClient
}

// SetHttpClient implements contracts.ClientConfig.
func (c *ClientConfigBase) SetHttpClient(httpc contracts.HTTPClient) {
	c.httpClient = httpc
}

// Validations implements contracts.ClientConfig.
func (c *ClientConfigBase) Validations() contracts.Validations {
	return c.validations
}

func newClientConfigBase(baseURL string) *ClientConfigBase {
	var validations []error

	if baseURL == "" {
		validations = append(validations, ErrEmptyBaseURL)
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		validations = append(validations, errors.Join(ErrParseURL, err))
	}

	return &ClientConfigBase{
		httpClient:    newDefaultHTTPClient(),
		httpHeader:    newDefaultHTTPHeader(),
		httpCookie:    newDefaultHttpCookies(),
		validations:   newDefaultValidations(validations),
		ConfigBaseURL: newDefaultBaseURL(parsedURL),
	}
}

func newBalancedClientConfigBase(baseURLs []string) *ClientConfigBase {
	var validations []error

	parsedURLs := make([]*url.URL, 0, len(baseURLs)) // pre-alloc cap like baseURLs

	for index, baseURL := range baseURLs {
		if baseURL == "" {
			validations = append(validations, fmt.Errorf("base URL %d: %w", index, ErrEmptyBaseURL))
			continue
		}

		parsedURL, err := url.Parse(baseURL)
		if err != nil {
			validations = append(validations, errors.Join(ErrParseURL, err))
		}

		parsedURLs = append(parsedURLs, parsedURL)
	}

	if len(parsedURLs) == 0 {
		validations = append(validations, ErrEmptyBaseURL)
	}

	return &ClientConfigBase{
		httpClient:    newDefaultHTTPClient(),
		httpHeader:    newDefaultHTTPHeader(),
		httpCookie:    newDefaultHttpCookies(),
		validations:   newDefaultValidations(validations),
		ConfigBaseURL: newBalancedBaseURL(parsedURLs),
	}
}
