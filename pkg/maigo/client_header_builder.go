package maigo

import (
	"github.com/jeanmolossi/maigo/pkg/maigo/contracts"
	"github.com/jeanmolossi/maigo/pkg/maigo/header"
	"github.com/jeanmolossi/maigo/pkg/maigo/mime"
)

var _ contracts.BuilderHeader[contracts.ClientBuilder] = (*ClientHeaderBuilder)(nil)

type ClientHeaderBuilder struct {
	parent *ClientBuilder
}

func (b *ClientBuilder) Header() contracts.BuilderHeader[contracts.ClientBuilder] {
	return &ClientHeaderBuilder{parent: b}
}

// Add implements contracts.BuilderHeader.
func (c *ClientHeaderBuilder) Add(key header.Type, value string) contracts.ClientBuilder {
	c.parent.client.Header().Add(key, value)
	return c.parent
}

// AddAccept implements contracts.BuilderHeader.
func (c *ClientHeaderBuilder) AddAccept(value mime.Type) contracts.ClientBuilder {
	c.Add(header.Accept, value.String())
	return c.parent
}

// AddAll implements contracts.BuilderHeader.
func (c *ClientHeaderBuilder) AddAll(headers map[header.Type]string) contracts.ClientBuilder {
	for key, value := range headers {
		c.Add(key, value)
	}

	return c.parent
}

// AddContentType implements contracts.BuilderHeader.
func (c *ClientHeaderBuilder) AddContentType(value mime.Type) contracts.ClientBuilder {
	c.Add(header.ContentType, value.String())
	return c.parent
}

// AddUserAgent implements contracts.BuilderHeader.
func (c *ClientHeaderBuilder) AddUserAgent(value string) contracts.ClientBuilder {
	c.Add(header.UserAgent, value)
	return c.parent
}

// Set implements contracts.BuilderHeader.
func (c *ClientHeaderBuilder) Set(key header.Type, value string) contracts.ClientBuilder {
	c.parent.client.Header().Set(key, value)
	return c.parent
}

// SetAll implements contracts.BuilderHeader.
func (c *ClientHeaderBuilder) SetAll(headers map[header.Type]string) contracts.ClientBuilder {
	for key, value := range headers {
		c.Set(key, value)
	}

	return c.parent
}
