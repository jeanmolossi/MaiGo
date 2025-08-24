package maigo

import (
	"context"
	"testing"
)

func TestContextSetIgnoresNil(t *testing.T) {
	t.Parallel()

	c := newDefaultContext()
	ctxBefore := c.Unwrap()

	c.Set(nil) //nolint:staticcheck // testing nil handling

	if c.Unwrap() != ctxBefore {
		t.Errorf("Set(nil) changed context: got %v, want %v", c.Unwrap(), ctxBefore)
	}
}

func TestContextSetNilReceiver(t *testing.T) {
	t.Parallel()

	var c *Context

	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("Set() panicked: %v", r)
			}
		}()
		c.Set(context.Background())
	}()
}

type ctxKey struct{}

func TestContextSetReplacesContext(t *testing.T) {
	t.Parallel()

	c := newDefaultContext()
	ctx := context.WithValue(context.Background(), ctxKey{}, "v")

	c.Set(ctx)

	if got := c.Unwrap(); got != ctx {
		t.Fatalf("Unwrap() = %v, want %v", got, ctx)
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

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Set(ctx)
	}
}

func BenchmarkContextUnwrap(b *testing.B) {
	c := newDefaultContext()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = c.Unwrap()
	}
}
