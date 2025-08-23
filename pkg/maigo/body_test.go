package maigo

import (
	"bytes"
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
}

func TestBufferedBodyJSONOperations(t *testing.T) {
	t.Parallel()

	body := newBufferedBody()
	in := sample{Name: "bob"}

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
}

func TestBufferedBodyXMLOperations(t *testing.T) {
	t.Parallel()

	body := newBufferedBody()
	in := sample{Name: "carol"}

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
	if err := body.WriteAsJSON(func() {}); err == nil {
		t.Error("WriteAsJSON() expected error for unsupported type")
	}
	// invalid XML write
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

	errCh := make(chan error, 50)

	for i := 0; i < 50; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			_, err := body.ReadAsString()
			errCh <- err
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			t.Fatalf("concurrent ReadAsString() error = %v", err)
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
}

func TestUnbufferedBodyXMLOperations(t *testing.T) {
	t.Parallel()

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

	errCh := make(chan error, 50)

	for i := 0; i < 50; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			_, err := body.ReadAsString()
			errCh <- err
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			t.Fatalf("concurrent ReadAsString() error = %v", err)
		}
	}
}
