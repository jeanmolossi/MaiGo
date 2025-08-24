package maigo

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"testing"
)

type sample struct {
	Name string `json:"name" xml:"name"`
}

func TestBufferedBodyStringOperations(t *testing.T) {
	t.Parallel()

	body := newBufferedBody()

	const want = "hello"

	if err := body.WriteAsString(want); err != nil {
		t.Fatalf("WriteAsString() error = %v", err)
	}

	got, err := body.ReadAsString()
	if err != nil {
		t.Fatalf("ReadAsString() error = %v", err)
	}

	if got != want {
		t.Errorf("ReadAsString() = %q, want %q", got, want)
	}

	// Body should not be consumed after read
	gotAgain, err := body.ReadAsString()
	if err != nil {
		t.Fatalf("ReadAsString() second read error = %v", err)
	}

	if gotAgain != want {
		t.Errorf("ReadAsString() second read = %q, want %q", gotAgain, want)
	}

	// Second write replaces previous content
	const newWant = "world"
	if err := body.WriteAsString(newWant); err != nil {
		t.Fatalf("WriteAsString() second write error = %v", err)
	}

	got, err = body.ReadAsString()
	if err != nil {
		t.Fatalf("ReadAsString() after rewrite error = %v", err)
	}

	if got != newWant {
		t.Errorf("ReadAsString() after rewrite = %q, want %q", got, newWant)
	}
}

func TestBufferedBodyJSONOperations(t *testing.T) {
	t.Parallel()

	t.Run("round trip", func(t *testing.T) {
		body := newBufferedBody()
		in := sample{Name: "bob"}

		if err := body.WriteAsJSON(in); err != nil {
			t.Fatalf("WriteAsJSON() error = %v", err)
		}

		var out sample
		if err := body.ReadAsJSON(&out); err != nil {
			t.Fatalf("ReadAsJSON() first read error = %v", err)
		}

		if out != in {
			t.Errorf("ReadAsJSON() first read = %#v, want %#v", out, in)
		}

		var outAgain sample
		if err := body.ReadAsJSON(&outAgain); err != nil {
			t.Fatalf("ReadAsJSON() second read error = %v", err)
		}

		if outAgain != in {
			t.Errorf("ReadAsJSON() second read = %#v, want %#v", outAgain, in)
		}
	})

	t.Run("rewrite replaces content", func(t *testing.T) {
		body := newBufferedBody()
		first := sample{Name: "first"}
		second := sample{Name: "second"}

		if err := body.WriteAsJSON(first); err != nil {
			t.Fatalf("WriteAsJSON() first write error = %v", err)
		}

		if err := body.WriteAsJSON(second); err != nil {
			t.Fatalf("WriteAsJSON() second write error = %v", err)
		}

		var out sample
		if err := body.ReadAsJSON(&out); err != nil {
			t.Fatalf("ReadAsJSON() after rewrite error = %v", err)
		}

		if out != second {
			t.Errorf("ReadAsJSON() after rewrite = %#v, want %#v", out, second)
		}
	})
}

func TestBufferedBodyXMLOperations(t *testing.T) {
	t.Parallel()

	t.Run("round trip", func(t *testing.T) {
		body := newBufferedBody()
		in := sample{Name: "carol"}

		if err := body.WriteAsXML(in); err != nil {
			t.Fatalf("WriteAsXML() error = %v", err)
		}

		var out sample
		if err := body.ReadAsXML(&out); err != nil {
			t.Fatalf("ReadAsXML() first read error = %v", err)
		}

		if out != in {
			t.Errorf("ReadAsXML() first read = %#v, want %#v", out, in)
		}

		var outAgain sample
		if err := body.ReadAsXML(&outAgain); err != nil {
			t.Fatalf("ReadAsXML() second read error = %v", err)
		}

		if outAgain != in {
			t.Errorf("ReadAsXML() second read = %#v, want %#v", outAgain, in)
		}
	})

	t.Run("rewrite replaces content", func(t *testing.T) {
		body := newBufferedBody()
		first := sample{Name: "first"}
		second := sample{Name: "second"}

		if err := body.WriteAsXML(first); err != nil {
			t.Fatalf("WriteAsXML() first write error = %v", err)
		}

		if err := body.WriteAsXML(second); err != nil {
			t.Fatalf("WriteAsXML() second write error = %v", err)
		}

		var out sample
		if err := body.ReadAsXML(&out); err != nil {
			t.Fatalf("ReadAsXML() after rewrite error = %v", err)
		}

		if out != second {
			t.Errorf("ReadAsXML() after rewrite = %#v, want %#v", out, second)
		}
	})
}

func TestBufferedBodyInvalidData(t *testing.T) {
	t.Parallel()

	body := newBufferedBody()
	// invalid JSON read
	if err := body.WriteAsString("not-json"); err != nil {
		t.Fatalf("WriteAsString() error = %v", err)
	}

	var js sample
	if err := body.ReadAsJSON(&js); err == nil {
		t.Error("ReadAsJSON() expected error, got nil")
	}

	// invalid XML read
	body = newBufferedBody()
	if err := body.WriteAsString("<invalid>"); err != nil {
		t.Fatalf("WriteAsString() error = %v", err)
	}

	var x sample
	if err := body.ReadAsXML(&x); err == nil {
		t.Error("ReadAsXML() expected error, got nil")
	}

	// invalid JSON write
	body = newBufferedBody()
	if err := body.WriteAsJSON(func() {}); err == nil {
		t.Error("WriteAsJSON() expected error for unsupported type")
	}

	// invalid XML write
	body = newBufferedBody()
	if err := body.WriteAsXML(func() {}); err == nil {
		t.Error("WriteAsXML() expected error for unsupported type")
	}
}

func TestBufferedBodyConcurrency(t *testing.T) {
	t.Parallel()

	body := newBufferedBody()

	const want = "concurrent"

	if err := body.WriteAsString(want); err != nil {
		t.Fatalf("WriteAsString() error = %v", err)
	}

	var wg sync.WaitGroup

	type result struct {
		s   string
		err error
	}

	resCh := make(chan result, 50)

	for i := 0; i < 50; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			s, err := body.ReadAsString()
			if err == nil && s != want {
				err = fmt.Errorf("ReadAsString() = %q, want %q", s, want)
			}

			resCh <- result{s: s, err: err}
		}()
	}

	wg.Wait()
	close(resCh)

	for res := range resCh {
		if res.err != nil {
			t.Fatalf("concurrent ReadAsString() error = %v", res.err)
		}

		if res.s != want {
			t.Errorf("ReadAsString() = %q, want %q", res.s, want)
		}
	}
}

func TestUnbufferedBodyStringOperations(t *testing.T) {
	t.Parallel()

	reader := io.NopCloser(bytes.NewBufferString(""))
	body := newUnbufferedBody(reader)

	const want = "hello"

	if err := body.WriteAsString(want); err != nil {
		t.Fatalf("WriteAsString() error = %v", err)
	}

	got, err := body.ReadAsString()
	if err != nil {
		t.Fatalf("ReadAsString() error = %v", err)
	}

	if got != want {
		t.Errorf("ReadAsString() = %q, want %q", got, want)
	}

	// Body should be reset after read
	gotAgain, err := body.ReadAsString()
	if err != nil {
		t.Fatalf("ReadAsString() second read error = %v", err)
	}

	if gotAgain != want {
		t.Errorf("ReadAsString() second read = %q, want %q", gotAgain, want)
	}
}

func TestUnbufferedBodyJSONOperations(t *testing.T) {
	t.Parallel()

	t.Run("round trip", func(t *testing.T) {
		body := newUnbufferedBody(io.NopCloser(bytes.NewBuffer(nil)))
		in := sample{Name: "dan"}

		if err := body.WriteAsJSON(in); err != nil {
			t.Fatalf("WriteAsJSON() error = %v", err)
		}

		var out sample
		if err := body.ReadAsJSON(&out); err != nil {
			t.Fatalf("ReadAsJSON() error = %v", err)
		}

		if out != in {
			t.Errorf("ReadAsJSON() = %#v, want %#v", out, in)
		}
	})

	t.Run("rewrite replaces content", func(t *testing.T) {
		body := newUnbufferedBody(io.NopCloser(bytes.NewBuffer(nil)))
		first := sample{Name: "first"}
		second := sample{Name: "second"}

		if err := body.WriteAsJSON(first); err != nil {
			t.Fatalf("WriteAsJSON() first write error = %v", err)
		}

		if err := body.WriteAsJSON(second); err != nil {
			t.Fatalf("WriteAsJSON() second write error = %v", err)
		}

		var out sample
		if err := body.ReadAsJSON(&out); err != nil {
			t.Fatalf("ReadAsJSON() after rewrite error = %v", err)
		}

		if out != second {
			t.Errorf("ReadAsJSON() after rewrite = %#v, want %#v", out, second)
		}
	})
}

func TestUnbufferedBodyXMLOperations(t *testing.T) {
	t.Parallel()

	t.Run("round trip", func(t *testing.T) {
		body := newUnbufferedBody(io.NopCloser(bytes.NewBuffer(nil)))
		in := sample{Name: "eric"}

		if err := body.WriteAsXML(in); err != nil {
			t.Fatalf("WriteAsXML() error = %v", err)
		}

		var out sample
		if err := body.ReadAsXML(&out); err != nil {
			t.Fatalf("ReadAsXML() error = %v", err)
		}

		if out != in {
			t.Errorf("ReadAsXML() = %#v, want %#v", out, in)
		}
	})

	t.Run("rewrite replaces content", func(t *testing.T) {
		body := newUnbufferedBody(io.NopCloser(bytes.NewBuffer(nil)))
		first := sample{Name: "first"}
		second := sample{Name: "second"}

		if err := body.WriteAsXML(first); err != nil {
			t.Fatalf("WriteAsXML() first write error = %v", err)
		}

		if err := body.WriteAsXML(second); err != nil {
			t.Fatalf("WriteAsXML() second write error = %v", err)
		}

		var out sample
		if err := body.ReadAsXML(&out); err != nil {
			t.Fatalf("ReadAsXML() after rewrite error = %v", err)
		}

		if out != second {
			t.Errorf("ReadAsXML() after rewrite = %#v, want %#v", out, second)
		}
	})
}

func TestUnbufferedBodyInvalidData(t *testing.T) {
	t.Parallel()

	body := newUnbufferedBody(io.NopCloser(bytes.NewBuffer(nil)))
	if err := body.WriteAsString("not-json"); err != nil {
		t.Fatalf("WriteAsString() error = %v", err)
	}

	var js sample
	if err := body.ReadAsJSON(&js); err == nil {
		t.Error("ReadAsJSON() expected error, got nil")
	}

	body = newUnbufferedBody(io.NopCloser(bytes.NewBuffer(nil)))
	if err := body.WriteAsString("<invalid>"); err != nil {
		t.Fatalf("WriteAsString() error = %v", err)
	}

	var x sample
	if err := body.ReadAsXML(&x); err == nil {
		t.Error("ReadAsXML() expected error, got nil")
	}

	if err := body.WriteAsJSON(func() {}); err == nil {
		t.Error("WriteAsJSON() expected error for unsupported type")
	}

	if err := body.WriteAsXML(func() {}); err == nil {
		t.Error("WriteAsXML() expected error for unsupported type")
	}
}

func TestUnbufferedBodyConcurrency(t *testing.T) {
	t.Parallel()

	body := newUnbufferedBody(io.NopCloser(bytes.NewBufferString("")))

	const want = "concurrent"

	if err := body.WriteAsString(want); err != nil {
		t.Fatalf("WriteAsString() error = %v", err)
	}

	var wg sync.WaitGroup

	type result struct {
		s   string
		err error
	}

	resCh := make(chan result, 50)

	for i := 0; i < 50; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			s, err := body.ReadAsString()
			if err == nil && s != want {
				err = fmt.Errorf("ReadAsString() = %q, want %q", s, want)
			}

			resCh <- result{s: s, err: err}
		}()
	}

	wg.Wait()
	close(resCh)

	for res := range resCh {
		if res.err != nil {
			t.Fatalf("concurrent ReadAsString() error = %v", res.err)
		}

		if res.s != want {
			t.Errorf("ReadAsString() = %q, want %q", res.s, want)
		}
	}
}
