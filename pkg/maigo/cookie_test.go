package maigo

import (
	"fmt"
	"net/http"
	"testing"
)

var (
	benchCookie  *http.Cookie
	benchCookies []*http.Cookie
)

func Test_newCookiesWithCapacity_Negative(t *testing.T) {
	t.Parallel()

	c := newCookiesWithCapacity(-1)
	if c == nil {
		t.Fatalf("newCookiesWithCapacity(-1) returned nil")
	}

	if c.Len() != 0 {
		t.Fatalf("Len() = %d, want %d", c.Len(), 0)
	}
}

func Test_isValidCookieName(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		valid bool
	}{
		{"session", true},
		{"token-123", true},
		{"", false},
		{"bad name", false},
		{";bad", false},
		{"bad;name", false},
	}

	for _, tc := range cases {
		if got := isValidCookieName(tc.name); got != tc.valid {
			t.Fatalf("isValidCookieName(%q) = %v, want %v", tc.name, got, tc.valid)
		}
	}
}

func TestCookies_AddAndCount(t *testing.T) {
	t.Parallel()

	c := newDefaultHTTPCookies()

	if c.Len() != 0 {
		t.Fatalf("Len() = %d, want %d", c.Len(), 0)
	}

	cookie := &http.Cookie{Name: "session", Value: "abc"}
	c.Add(cookie)

	if c.Len() != 1 {
		t.Fatalf("Len() = %d, want %d", c.Len(), 1)
	}

	// adding a nil cookie should not change count
	c.Add(nil)

	// adding a cookie with empty name should not change count
	c.Add(&http.Cookie{Value: "no-name"})

	// adding a cookie with whitespace-only name should not change count
	c.Add(&http.Cookie{Name: "   ", Value: "blank-name"})

	// adding a cookie with invalid characters should not change count
	c.Add(&http.Cookie{Name: "bad;name", Value: "x"})

	if c.Len() != 1 {
		t.Fatalf("after invalid Add Len() = %d, want %d", c.Len(), 1)
	}

	if c.Len() != c.Count() {
		t.Fatalf("Len() = %d, Count() = %d", c.Len(), c.Count())
	}
}

func TestCookies_Add_PreservesName(t *testing.T) {
	t.Parallel()

	c := newDefaultHTTPCookies()
	orig := &http.Cookie{Name: "  session  ", Value: "v"}

	c.Add(orig)

	if orig.Name != "  session  " {
		t.Fatalf("Add mutated original Name = %q, want %q", orig.Name, "  session  ")
	}

	if c.Len() != 1 {
		t.Fatalf("Len() after Add = %d, want %d", c.Len(), 1)
	}

	got := c.Get(0)
	if got == nil {
		t.Fatalf("Get(0) returned nil")
	}

	if got.Name != orig.Name {
		t.Fatalf("stored Name = %q, want %q", got.Name, orig.Name)
	}
}

func TestCookies_Add_InvalidName(t *testing.T) {
	t.Parallel()

	c := newDefaultHTTPCookies()

	invalid := []*http.Cookie{
		{Name: "bad name", Value: "v"},
		{Name: "bad;name", Value: "v"},
		{Name: "", Value: "v"},
	}

	for _, ck := range invalid {
		c.Add(ck)
	}

	if c.Len() != 0 {
		t.Fatalf("Len() = %d, want %d", c.Len(), 0)
	}
}

func TestCookies_Add_ClonesInput(t *testing.T) {
	t.Parallel()

	c := newDefaultHTTPCookies()
	ck := &http.Cookie{Name: "session", Value: "abc"}
	c.Add(ck)

	ck.Value = "changed"

	got := c.Get(0)
	if got == nil || got.Value != "abc" {
		var v string
		if got != nil {
			v = got.Value
		}

		t.Fatalf("stored cookie Value = %q, want %q", v, "abc")
	}
}

func TestCookies_UnwrapEmpty(t *testing.T) {
	t.Parallel()

	c := newDefaultHTTPCookies()

	if got := c.Unwrap(); got != nil {
		t.Fatalf("Unwrap() = %v, want nil", got)
	}
}

func TestCookies_NilReceiver(t *testing.T) {
	t.Parallel()

	var c *Cookies

	c.Add(&http.Cookie{Name: "k", Value: "v"})

	if c.Count() != 0 {
		t.Fatalf("Count() on nil receiver = %d, want %d", c.Count(), 0)
	}

	if c.Len() != 0 {
		t.Fatalf("Len() on nil receiver = %d, want %d", c.Len(), 0)
	}

	if c.Get(0) != nil {
		t.Fatalf("Get on nil receiver returned non-nil")
	}

	if c.Unwrap() != nil {
		t.Fatalf("Unwrap on nil receiver returned non-nil slice")
	}
}

func TestCookies_Get(t *testing.T) {
	t.Parallel()

	c := newDefaultHTTPCookies()
	cookie := &http.Cookie{Name: "session", Value: "abc"}
	c.Add(cookie)

	got := c.Get(0)
	if got == nil || got.Name != cookie.Name || got.Value != cookie.Value {
		var gotName, gotValue string
		if got != nil {
			gotName, gotValue = got.Name, got.Value
		}

		t.Fatalf("Get(0) = {Name: %q, Value: %q}, want {Name: %q, Value: %q}", gotName, gotValue, cookie.Name, cookie.Value)
	}

	if got := c.Get(1); got != nil {
		t.Fatalf("Get(1) = %v, want nil", got)
	}

	if got := c.Get(-1); got != nil {
		t.Fatalf("Get(-1) = %v, want nil", got)
	}
}

func TestCookies_Get_ReturnsCopy(t *testing.T) {
	t.Parallel()

	c := newDefaultHTTPCookies()
	c.Add(&http.Cookie{Name: "session", Value: "abc"})

	got := c.Get(0)
	if got == nil {
		t.Fatalf("Get(0) returned nil")
	}

	got.Value = "changed"

	again := c.Get(0)
	if again == nil || again.Value != "abc" {
		var v string
		if again != nil {
			v = again.Value
		}

		t.Fatalf("after mutating copy, Get(0).Value = %q, want %q", v, "abc")
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

	unwrapped = append(unwrapped, &http.Cookie{Name: "extra", Value: "v"})
	if len(unwrapped) != 2 {
		t.Fatalf("appended slice length = %d, want %d", len(unwrapped), 2)
	}

	got := c.Get(0)
	if got.Name != "session" || got.Unparsed[0] != "a=b" {
		t.Fatalf("original cookie mutated: %v", got)
	}

	if c.Len() != 1 {
		t.Fatalf("original slice length changed: %d", c.Len())
	}

	if c.Get(1) != nil {
		t.Fatalf("appending to unwrapped slice affected original: %v", c.Get(1))
	}
}

func TestCookies_Unwrap_PreservesEmptyUnparsed(t *testing.T) {
	t.Parallel()

	c := newDefaultHTTPCookies()
	c.Add(&http.Cookie{Name: "n", Value: "v", Unparsed: []string{}})

	u := c.Unwrap()
	if len(u) != 1 {
		t.Fatalf("Unwrap length = %d, want %d", len(u), 1)
	}

	if u[0].Unparsed == nil {
		t.Fatalf("Unwrap returned nil Unparsed slice")
	}

	if len(u[0].Unparsed) != 0 {
		t.Fatalf("Unwrap returned Unparsed len = %d, want %d", len(u[0].Unparsed), 0)
	}
}

func BenchmarkCookies_Add(b *testing.B) {
	cookie := &http.Cookie{Name: "k", Value: "v"}

	b.Run("prealloc", func(b *testing.B) {
		c := newCookiesWithCapacity(b.N)

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			c.Add(cookie)
		}
	})

	b.Run("growth", func(b *testing.B) {
		c := &Cookies{}

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			c.Add(cookie)
		}
	})
}

func BenchmarkCookies_Get(b *testing.B) {
	c := newDefaultHTTPCookies()
	c.Add(&http.Cookie{Name: "k", Value: "v"})
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchCookie = c.Get(0)
	}
}

func BenchmarkCookies_Unwrap(b *testing.B) {
	sizes := []int{0, 1, 8, 64, 512}

	for _, n := range sizes {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			c := newDefaultHTTPCookies()
			for i := 0; i < n; i++ {
				c.Add(&http.Cookie{Name: "k", Value: "v"})
			}

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				benchCookies = c.Unwrap()
			}
		})
	}
}
