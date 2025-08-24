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

	// adding a cookie with empty name should not change count
	c.Add(&http.Cookie{Value: "no-name"})

	if c.Count() != 1 {
		t.Fatalf("after invalid Add Count() = %d, want %d", c.Count(), 1)
	}
}

func TestCookies_UnwrapEmpty(t *testing.T) {
	t.Parallel()

	c := newDefaultHTTPCookies()

	if got := c.Unwrap(); got != nil {
		t.Fatalf("Unwrap() = %v, want nil", got)
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

func TestCookies_UnwrapDeepCopy(t *testing.T) {
	t.Parallel()

	c := newDefaultHTTPCookies()
	cookie := &http.Cookie{Name: "session", Value: "abc", Unparsed: []string{"a=b"}}
	c.Add(cookie)

	unwrapped := c.Unwrap()
	if len(unwrapped) != 1 {
		t.Fatalf("Unwrap length = %d, want %d", len(unwrapped), 1)
	}

	if unwrapped[0] == cookie {
		t.Fatalf("Unwrap returned original pointer")
	}

	unwrapped[0].Name = "changed"
	unwrapped[0].Unparsed[0] = "x=y"

	got := c.Get(0)
	if got.Name != "session" || got.Unparsed[0] != "a=b" {
		t.Fatalf("original cookie mutated: %v", got)
	}
}

func BenchmarkCookies_Add(b *testing.B) {
	c := newDefaultHTTPCookies()
	cookie := &http.Cookie{Name: "k", Value: "v"}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.cookies = c.cookies[:0]
		c.Add(cookie)
	}
}

func BenchmarkCookies_Get(b *testing.B) {
	c := newDefaultHTTPCookies()
	c.Add(&http.Cookie{Name: "k", Value: "v"})
	b.ReportAllocs()
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

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = c.Unwrap()
	}
}
