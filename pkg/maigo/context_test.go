package maigo

import (
	"context"
	"testing"
)

func TestContextSetIgnoresNil(t *testing.T) {
	t.Parallel()

	c := newDefaultContext()
	ctxBefore := c.Unwrap()

	c.Set(nil)

	if c.Unwrap() != ctxBefore {
		t.Errorf("Set(nil) changed context: got %v, want %v", c.Unwrap(), ctxBefore)
	}
}

func TestContextUnwrapReturnsNonNil(t *testing.T) {
	t.Parallel()

	var c Context

	if c.Unwrap() == nil {
		t.Error("Unwrap() returned nil")
	}

	var cp *Context
	if cp.Unwrap() == nil {
		t.Error("Unwrap() on nil receiver returned nil")
	}
}

func BenchmarkContextSet(b *testing.B) {
	c := newDefaultContext()
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		c.Set(ctx)
	}
}

func BenchmarkContextUnwrap(b *testing.B) {
	c := newDefaultContext()

	for i := 0; i < b.N; i++ {
		_ = c.Unwrap()
	}
}
