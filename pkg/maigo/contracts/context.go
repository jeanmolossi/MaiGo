package contracts

import "context"

// Context is the interface that wraps the basic methods to manage HTTP request context.
//
// This wrapper provides a layer of abstraction over the standard context.Context,
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
	// Unwrap exposes the underlying context.
	Unwrap() context.Context
	// Set replaces the underlying context.
	Set(ctx context.Context)
}
