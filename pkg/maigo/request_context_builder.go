package maigo

import (
	"context"

	"github.com/jeanmolossi/maigo/pkg/maigo/contracts"
)

var _ contracts.BuilderRequestContext[contracts.RequestBuilder] = (*RequestContextBuilder)(nil)

type RequestContextBuilder struct {
	parent *RequestBuilder
	config *RequestConfigBase
}

func (r *RequestBuilder) Context() contracts.BuilderRequestContext[contracts.RequestBuilder] {
	return &RequestContextBuilder{
		parent: r,
		config: r.request.config,
	}
}

// Set implements contracts.BuilderRequestContext.
func (r *RequestContextBuilder) Set(ctx context.Context) contracts.RequestBuilder {
	r.config.Context().Set(ctx)
	return r.parent
}
