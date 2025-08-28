package maigo

import (
	"sync"
	"testing"

	"github.com/jeanmolossi/maigo/pkg/maigo/header"
)

func TestHeader_ConcurrentAddGet(t *testing.T) {
	t.Parallel()

	h := newDefaultHTTPHeader()

	const (
		goroutines = 8
		iterations = 1000
	)

	var wg sync.WaitGroup

	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()

			for j := 0; j < iterations; j++ {
				h.Add(header.ContentType, "application/json")
				_ = h.Get(header.ContentType)
			}
		}()
	}

	wg.Wait()
}

func TestHeader_AddAndGet(t *testing.T) {
	t.Parallel()

	h := newDefaultHTTPHeader()

	h.Add(header.ContentType, "application/json")
	h.Add(header.ContentType, "text/html")

	if got := h.Get(header.ContentType); got != "application/json" {
		t.Errorf("Get() = %q, want %q", got, "application/json")
	}

	values := h.Unwrap().Values(header.ContentType.String())
	if len(values) != 2 {
		t.Fatalf("values length = %d, want 2", len(values))
	}

	if values[0] != "application/json" || values[1] != "text/html" {
		t.Errorf("values = %#v, want [application/json text/html]", values)
	}
}

func TestHeader_Set(t *testing.T) {
	t.Parallel()

	h := newDefaultHTTPHeader()

	h.Add(header.ContentType, "application/json")
	h.Set(header.ContentType, "text/plain")

	if got := h.Get(header.ContentType); got != "text/plain" {
		t.Errorf("Get() after Set = %q, want %q", got, "text/plain")
	}

	values := h.Unwrap().Values(header.ContentType.String())
	if len(values) != 1 || values[0] != "text/plain" {
		t.Errorf("values = %#v, want [text/plain]", values)
	}
}

func TestHeader_Get_NotFound(t *testing.T) {
	t.Parallel()

	h := newDefaultHTTPHeader()

	if got := h.Get(header.ContentType); got != "" {
		t.Errorf("Get() = %q, want empty", got)
	}
}

func TestHeader_UnwrapReturnsCopy(t *testing.T) {
	t.Parallel()

	h := newDefaultHTTPHeader()
	h.Set(header.ContentType, "application/json")

	unwrapped := h.Unwrap()
	if unwrapped == nil {
		t.Fatal("Unwrap() returned nil")
	}

	if got := unwrapped.Get(header.ContentType.String()); got != "application/json" {
		t.Fatalf("unwrapped header value = %q, want %q", got, "application/json")
	}

	unwrapped.Set(header.ContentType.String(), "text/plain")

	if got := h.Get(header.ContentType); got != "application/json" {
		t.Errorf("original header modified after Unwrap(), got %q", got)
	}
}

func TestHeader_NilReceiver(t *testing.T) {
	t.Parallel()

	var h *Header

	h.Add(header.ContentType, "application/json")
	h.Set(header.ContentType, "text/plain")

	if got := h.Get(header.ContentType); got != "" {
		t.Errorf("Get() on nil receiver = %q, want empty", got)
	}
}

func TestHeader_UnwrapNilReceiver(t *testing.T) {
	t.Parallel()

	var h *Header

	hdr := h.Unwrap()
	if hdr == nil {
		t.Fatal("Unwrap() on nil receiver returned nil")
	}

	if len(*hdr) != 0 {
		t.Fatalf("Unwrap() on nil receiver returned non-empty header: %#v", *hdr)
	}
}

func TestHeader_InvalidInputIgnored(t *testing.T) {
	t.Parallel()

	h := newDefaultHTTPHeader()
	badName := header.Parse("Bad\nName")
	h.Set(badName, "value")

	if v := h.Get(badName); v != "" {
		t.Fatalf("expected no value for invalid name, got %q", v)
	}

	h.Set(header.ContentType, "text/plain\r\nmalicious")

	if v := h.Get(header.ContentType); v != "" {
		t.Fatalf("expected no value for invalid value, got %q", v)
	}
}

func BenchmarkHeaderAdd(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		h := newDefaultHTTPHeader()
		h.Add(header.ContentType, "application/json")
	}
}

func BenchmarkHeaderSet(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		h := newDefaultHTTPHeader()
		h.Set(header.ContentType, "application/json")
	}
}

func BenchmarkHeaderAddReuse(b *testing.B) {
	b.ReportAllocs()

	h := newDefaultHTTPHeader()

	const batch = 1000

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		h.Add(header.ContentType, "application/json")

		if i%batch == batch-1 {
			b.StopTimer()
			h.Set(header.ContentType, "")
			b.StartTimer()
		}
	}

	b.StopTimer()
	h.Set(header.ContentType, "")
}

func BenchmarkHeaderSetReuse(b *testing.B) {
	b.ReportAllocs()

	h := newDefaultHTTPHeader()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		h.Set(header.ContentType, "application/json")
	}
}

func BenchmarkHeaderGet(b *testing.B) {
	b.ReportAllocs()

	h := newDefaultHTTPHeader()
	h.Set(header.ContentType, "application/json")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if h.Get(header.ContentType) == "" {
			b.Fatal("unexpected empty value")
		}
	}
}

func BenchmarkHeaderUnwrap(b *testing.B) {
	b.ReportAllocs()

	h := newDefaultHTTPHeader()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if h.Unwrap() == nil {
			b.Fatal("unwrap returned nil")
		}
	}
}
