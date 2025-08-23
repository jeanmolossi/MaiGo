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
	if ctx != nil {
		c.ctx = ctx
	}
}

// Unwrap implements contracts.Context.
func (c *Context) Unwrap() context.Context {
	return c.ctx
}

func newDefaultContext() *Context {
	return &Context{
		ctx: context.Background(),
	}
}
