// Package maigo contains core primitives for the MaiGo project.
//
// The Context type wraps a standard context.Context. It is intentionally
// nil-safe: calling Set(nil) leaves the existing context unchanged and Unwrap
// always yields a non-nil context. The type is not safe for concurrent use; the
// caller must not invoke Set and Unwrap concurrently.
package maigo

import (
	"context"
	"reflect"

	"github.com/jeanmolossi/maigo/pkg/maigo/contracts"
)

var _ contracts.Context = (*Context)(nil)

type Context struct {
	ctx context.Context
}

// Set replaces the wrapped context.
//
// Passing nil is a no-op. Context is not safe for concurrent use; callers must
// not invoke Set and Unwrap from multiple goroutines without synchronization.
func (c *Context) Set(ctx context.Context) {
	if c == nil || ctx == nil {
		return
	}

	if c.ctx != nil {
		oldVal := reflect.ValueOf(c.ctx)
		newVal := reflect.ValueOf(ctx)

		if oldVal.IsValid() && newVal.IsValid() &&
			oldVal.Kind() == reflect.Pointer && newVal.Kind() == reflect.Pointer &&
			oldVal.Pointer() == newVal.Pointer() {
			return
		}
	}

	c.ctx = ctx
}

// Unwrap returns the stored context or context.Background if the receiver or
// its stored context is nil. Returning Background ensures callers always receive
// a usable context and repeated calls in this nil case yield the same instance,
// providing idempotency. Cancellation and deadlines are expected to be managed
// by the caller.
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
