package maigo

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

type sample struct {
	Name string `json:"name" xml:"name"`
}

type trackedRC struct {
	io.Reader
	closed int32
}

func (t *trackedRC) Close() error {
	atomic.StoreInt32(&t.closed, 1)
	return nil
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

	t.Run("replaces reader closes previous", func(t *testing.T) {
		t.Parallel()

		prev := &trackedRC{Reader: strings.NewReader("old")}
		body := newUnbufferedBody(prev)

		if err := body.WriteAsString("new"); err != nil {
			t.Fatalf("WriteAsString() error = %v", err)
		}

		if atomic.LoadInt32(&prev.closed) != 1 {
			t.Error("WriteAsString() did not close previous reader")
		}

		prev = &trackedRC{Reader: strings.NewReader("old")}
		body = newUnbufferedBody(prev)

		if err := body.Set(strings.NewReader("new")); err != nil {
			t.Fatalf("Set() error = %v", err)
		}

		if atomic.LoadInt32(&prev.closed) != 1 {
			t.Error("Set() did not close previous reader")
		}
	})
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

	t.Run("write closes previous reader", func(t *testing.T) {
		t.Parallel()

		prev := &trackedRC{Reader: strings.NewReader("old")}
		body := newUnbufferedBody(prev)

		if err := body.WriteAsJSON(sample{Name: "new"}); err != nil {
			t.Fatalf("WriteAsJSON() error = %v", err)
		}

		if atomic.LoadInt32(&prev.closed) != 1 {
			t.Error("WriteAsJSON() did not close previous reader")
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

	t.Run("write closes previous reader", func(t *testing.T) {
		t.Parallel()

		prev := &trackedRC{Reader: strings.NewReader("old")}
		body := newUnbufferedBody(prev)

		if err := body.WriteAsXML(sample{Name: "new"}); err != nil {
			t.Fatalf("WriteAsXML() error = %v", err)
		}

		if atomic.LoadInt32(&prev.closed) != 1 {
			t.Error("WriteAsXML() did not close previous reader")
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

	if s, err := body.ReadAsString(); err != nil {
		t.Fatalf("ReadAsString() after JSON error = %v", err)
	} else if s != "not-json" {
		t.Errorf("ReadAsString() after JSON error = %q, want %q", s, "not-json")
	}

	body = newUnbufferedBody(io.NopCloser(bytes.NewBuffer(nil)))
	if err := body.WriteAsString("<invalid>"); err != nil {
		t.Fatalf("WriteAsString() error = %v", err)
	}

	var x sample
	if err := body.ReadAsXML(&x); err == nil {
		t.Error("ReadAsXML() expected error, got nil")
	}

	if s, err := body.ReadAsString(); err != nil {
		t.Fatalf("ReadAsString() after XML error = %v", err)
	} else if s != "<invalid>" {
		t.Errorf("ReadAsString() after XML error = %q, want %q", s, "<invalid>")
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
