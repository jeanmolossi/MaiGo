package maigo

import (
	"github.com/jeanmolossi/MaiGo/pkg/maigo/contracts"
	"github.com/jeanmolossi/MaiGo/pkg/maigo/header"
	"github.com/jeanmolossi/MaiGo/pkg/maigo/mime"
)

var _ contracts.BuilderHeader[contracts.RequestBuilder] = (*RequestHeaderBuilder)(nil)

type RequestHeaderBuilder struct {
	parent *RequestBuilder
	config *RequestConfigBase
}

func (r *RequestBuilder) Header() contracts.BuilderHeader[contracts.RequestBuilder] {
	return &RequestHeaderBuilder{
		parent: r,
		config: r.request.config,
	}
}

// Add implements contracts.BuilderHeader.
func (r *RequestHeaderBuilder) Add(key header.Type, value string) contracts.RequestBuilder {
	r.config.httpHeader.Add(key, value)
	return r.parent
}

// AddAccept implements contracts.BuilderHeader.
func (r *RequestHeaderBuilder) AddAccept(value mime.Type) contracts.RequestBuilder {
	r.config.httpHeader.Add(header.Accept, value.String())
	return r.parent
}

// AddAll implements contracts.BuilderHeader.
func (r *RequestHeaderBuilder) AddAll(headers map[header.Type]string) contracts.RequestBuilder {
	for key, value := range headers {
		r.Add(key, value)
	}

	return r.parent
}

// AddContentType implements contracts.BuilderHeader.
func (r *RequestHeaderBuilder) AddContentType(value mime.Type) contracts.RequestBuilder {
	r.Add(header.ContentType, value.String())
	return r.parent
}

// AddUserAgent implements contracts.BuilderHeader.
func (r *RequestHeaderBuilder) AddUserAgent(value string) contracts.RequestBuilder {
	r.Add(header.UserAgent, value)
	return r.parent
}

// Set implements contracts.BuilderHeader.
func (r *RequestHeaderBuilder) Set(key header.Type, value string) contracts.RequestBuilder {
	r.config.httpHeader.Set(key, value)
	return r.parent
}

// SetAll implements contracts.BuilderHeader.
func (r *RequestHeaderBuilder) SetAll(headers map[header.Type]string) contracts.RequestBuilder {
	for key, value := range headers {
		r.Set(key, value)
	}

	return r.parent
}
