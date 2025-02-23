package client

import (
	"errors"
	"net/url"

	"github.com/jeanmolossi/MaiGo/pkg/client/contracts"
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
	panic("unimplemented")
}

// DELETE implements contracts.ClientHTTPMethods.
func (c *ClientConfigBase) DELETE(path string) contracts.RequestBuilder {
	panic("unimplemented")
}

// GET implements contracts.ClientHTTPMethods.
func (c *ClientConfigBase) GET(path string) contracts.RequestBuilder {
	panic("unimplemented")
}

// HEAD implements contracts.ClientHTTPMethods.
func (c *ClientConfigBase) HEAD(path string) contracts.RequestBuilder {
	panic("unimplemented")
}

// OPTIONS implements contracts.ClientHTTPMethods.
func (c *ClientConfigBase) OPTIONS(path string) contracts.RequestBuilder {
	panic("unimplemented")
}

// PATCH implements contracts.ClientHTTPMethods.
func (c *ClientConfigBase) PATCH(path string) contracts.RequestBuilder {
	panic("unimplemented")
}

// POST implements contracts.ClientHTTPMethods.
func (c *ClientConfigBase) POST(path string) contracts.RequestBuilder {
	panic("unimplemented")
}

// PUT implements contracts.ClientHTTPMethods.
func (c *ClientConfigBase) PUT(path string) contracts.RequestBuilder {
	panic("unimplemented")
}

// TRACE implements contracts.ClientHTTPMethods.
func (c *ClientConfigBase) TRACE(path string) contracts.RequestBuilder {
	panic("unimplemented")
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
