package maigo

import (
	"net/http"
	"testing"
)

func TestCookies_AddAndCount(t *testing.T) {
	t.Parallel()

	c := newDefaultHTTPCookies()

	if c.Count() != 0 {
		t.Fatalf("Count() = %d, want %d", c.Count(), 0)
	}

	cookie := &http.Cookie{Name: "session", Value: "abc"}
	c.Add(cookie)

	if c.Count() != 1 {
		t.Fatalf("Count() = %d, want %d", c.Count(), 1)
	}

	// adding a nil cookie should not change count
	c.Add(nil)

	if c.Count() != 1 {
		t.Fatalf("after nil Add Count() = %d, want %d", c.Count(), 1)
	}
}

func TestCookies_Get(t *testing.T) {
	t.Parallel()

	c := newDefaultHTTPCookies()
	cookie := &http.Cookie{Name: "session", Value: "abc"}
	c.Add(cookie)

	if got := c.Get(0); got != cookie {
		t.Fatalf("Get(0) = %v, want %v", got, cookie)
	}

	if got := c.Get(1); got != nil {
		t.Fatalf("Get(1) = %v, want nil", got)
	}

	if got := c.Get(-1); got != nil {
		t.Fatalf("Get(-1) = %v, want nil", got)
	}
}

func TestCookies_Unwrap(t *testing.T) {
	t.Parallel()

	c := newDefaultHTTPCookies()
	cookie := &http.Cookie{Name: "session", Value: "abc"}
	c.Add(cookie)

	unwrapped := c.Unwrap()
	if len(unwrapped) != 1 {
		t.Fatalf("Unwrap length = %d, want %d", len(unwrapped), 1)
	}

	// ensure modifications to unwrapped slice do not affect original
	unwrapped[0] = &http.Cookie{Name: "changed", Value: "xyz"}

	got := c.Get(0)
	if got.Name != "session" || got.Value != "abc" {
		t.Fatalf("internal cookie modified: %v", got)
	}
}

func BenchmarkCookies_Add(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := newDefaultHTTPCookies()
		c.Add(&http.Cookie{Name: "k", Value: "v"})
	}
}

func BenchmarkCookies_Get(b *testing.B) {
	c := newDefaultHTTPCookies()
	c.Add(&http.Cookie{Name: "k", Value: "v"})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = c.Get(0)
	}
}

func BenchmarkCookies_Unwrap(b *testing.B) {
	c := newDefaultHTTPCookies()
	for i := 0; i < 10; i++ {
		c.Add(&http.Cookie{Name: "k", Value: "v"})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = c.Unwrap()
	}
}
