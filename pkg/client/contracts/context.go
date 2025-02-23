package contracts

import "context"

// Context is the interface that wraps the basic methods to manage HTTP request context.
//
// This wrapper providers a layer of abstraction over the stantard context.Context,
// allowing for potential enhancements to context handling without affecting the public API.
//
// It enables the package to implement custom context-related features, such as automatic context propagation
// or context-based tracing, while maintaining a simple interface.
//
// Example:
//
//	type TracingContext struct {
//	    ctx context.Context
//	    tracer Tracer
//	}
//
//	func (tc *TracingContext) Unwrap() context.Context {
//	    return tc.tracer.ContextWithSpan(tc.ctx)
//	}
type Context interface {
	Unwrap() context.Context
	Set(ctx context.Context)
}
