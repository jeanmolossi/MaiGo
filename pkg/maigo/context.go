package maigo

import (
	"context"

	"github.com/jeanmolossi/maigo/pkg/maigo/contracts"
)

var _ contracts.Context = (*Context)(nil)

type Context struct {
	ctx context.Context
}

// Set implements contracts.Context.
func (c *Context) Set(ctx context.Context) {
	if c == nil || ctx == nil {
		return
	}

	c.ctx = ctx
}

// Unwrap implements contracts.Context.
func (c *Context) Unwrap() context.Context {
	if c == nil || c.ctx == nil {
		return context.Background()
	}

	return c.ctx
}

func newDefaultContext() *Context {
	return &Context{
		ctx: context.Background(),
	}
}
