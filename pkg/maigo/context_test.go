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

type uncomparableCtx struct {
	context.Context
	v []int
}

func TestContextSetWithUncomparableContext(t *testing.T) {
	t.Parallel()

	var c Context

	ctx1 := uncomparableCtx{Context: context.Background(), v: []int{1}}
	c.Set(ctx1)

	ctx2 := uncomparableCtx{Context: context.Background(), v: []int{2}}

	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("Set() panicked: %v", r)
			}
		}()

		c.Set(ctx2)
	}()

	got := c.Unwrap().(uncomparableCtx)
	if got.v[0] != 2 {
		t.Fatalf("Unwrap() = %v, want %v", got.v, ctx2.v)
	}
}

func TestContextUnwrapReturnsNonNil(t *testing.T) {
	t.Parallel()

	var c Context

	first := c.Unwrap()
	if first == nil {
		t.Fatal("Unwrap() returned nil")
	}

	second := c.Unwrap()
	if second == nil {
		t.Fatal("Unwrap() returned nil on second call")
	}

	if first != second {
		t.Errorf("Unwrap() returned different instances: %p != %p", first, second)
	}

	var cp *Context

	firstNil := cp.Unwrap()
	if firstNil == nil {
		t.Fatal("Unwrap() on nil receiver returned nil")
	}

	secondNil := cp.Unwrap()
	if secondNil == nil {
		t.Fatal("Unwrap() on nil receiver returned nil on second call")
	}

	if firstNil != secondNil {
		t.Errorf("Unwrap() on nil receiver returned different instances: %p != %p", firstNil, secondNil)
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
