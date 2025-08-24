package maigo

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/jeanmolossi/maigo/pkg/maigo/contracts"
)

var (
	_ contracts.Body = (*BufferedBody)(nil)
	_ contracts.Body = (*UnbufferedBody)(nil)
)

type (
	BufferedBody struct {
		buffer *bytes.Buffer
		mutex  sync.RWMutex
	}

	UnbufferedBody struct {
		reader io.ReadCloser
		mutex  sync.RWMutex
	}
)

// Close implements contracts.Body.
func (u *UnbufferedBody) Close() (err error) {
	u.mutex.Lock()
	err = u.reader.Close()
	u.mutex.Unlock()

	return
}

// Read implements contracts.Body.
func (u *UnbufferedBody) Read(p []byte) (n int, err error) {
	u.mutex.Lock()
	n, err = u.reader.Read(p)
	u.mutex.Unlock()

	return
}

// ReadAsJSON reads the full body, decodes it into obj and replaces the
// underlying reader so the data can be consumed again. Decoding directly from
// u.reader would advance it, making subsequent reads return no data.
func (u *UnbufferedBody) ReadAsJSON(obj any) (err error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	prev := u.reader

	data, err := io.ReadAll(prev)
	if err != nil {
		return fmt.Errorf("failed reading body as JSON: %w", err)
	}

	if err = json.Unmarshal(data, obj); err != nil {
		return fmt.Errorf("failed decoding body as JSON: %w", err)
	}

	u.reader = io.NopCloser(bytes.NewReader(data))

	_ = prev.Close()

	return nil
}

// WriteAsJSON implements contracts.Body.
func (u *UnbufferedBody) WriteAsJSON(obj any) (err error) {
	var buf bytes.Buffer

	u.mutex.Lock()
	defer u.mutex.Unlock()

	err = json.NewEncoder(&buf).Encode(obj)
	if err != nil {
		return
	}

	u.reader = io.NopCloser(&buf)

	return
}

// ReadAsXML reads the full body, decodes it into obj and replaces the reader so
// future reads see the same data. Without copying the bytes, decoding would
// consume u.reader.
func (u *UnbufferedBody) ReadAsXML(obj any) (err error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	prev := u.reader

	data, err := io.ReadAll(prev)
	if err != nil {
		return fmt.Errorf("failed reading body as XML: %w", err)
	}

	if err = xml.Unmarshal(data, obj); err != nil {
		return fmt.Errorf("failed decoding body as XML: %w", err)
	}

	u.reader = io.NopCloser(bytes.NewReader(data))

	_ = prev.Close()

	return nil
}

// WriteAsXML implements contracts.Body.
func (u *UnbufferedBody) WriteAsXML(obj any) (err error) {
	var buf bytes.Buffer

	u.mutex.Lock()
	defer u.mutex.Unlock()

	err = xml.NewEncoder(&buf).Encode(obj)
	if err != nil {
		return
	}

	u.reader = io.NopCloser(&buf)

	return
}

// ReadAsString implements contracts.Body.
func (u *UnbufferedBody) ReadAsString() (string, error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	stringBytes, err := io.ReadAll(u.reader)
	if err != nil {
		return "", fmt.Errorf("failed reading body as string: %w", err)
	}

	u.reader = io.NopCloser(bytes.NewReader(stringBytes))

	return string(stringBytes), nil
}

// WriteAsString implements contracts.Body.
func (u *UnbufferedBody) WriteAsString(body string) (err error) {
	u.mutex.Lock()
	u.reader = io.NopCloser(strings.NewReader(body))
	u.mutex.Unlock()

	return
}

// Set implements contracts.Body.
func (u *UnbufferedBody) Set(body io.Reader) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if closer, ok := body.(io.ReadCloser); ok {
		u.reader = closer
	} else {
		u.reader = io.NopCloser(body)
	}

	return nil
}

// Unwrap implements contracts.Body.
func (u *UnbufferedBody) Unwrap() io.Reader {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	return u.reader
}

func newUnbufferedBody(reader io.ReadCloser) *UnbufferedBody {
	return &UnbufferedBody{
		reader: reader,
	}
}

// -----------------------------------------------------
//
// BufferedBody methods
//
// -----------------------------------------------------

// Close implements contracts.Body.
func (b *BufferedBody) Close() error {
	b.buffer.Reset()

	return nil
}

// Read implements contracts.Body.
func (b *BufferedBody) Read(p []byte) (n int, err error) {
	b.mutex.RLock()
	n, err = b.buffer.Read(p)
	b.mutex.RUnlock()

	return
}

// ReadAsJSON implements contracts.Body.
func (b *BufferedBody) ReadAsJSON(obj any) (err error) {
	b.mutex.RLock()
	err = json.Unmarshal(b.buffer.Bytes(), obj)
	b.mutex.RUnlock()

	return
}

// WriteAsJSON implements contracts.Body.
func (b *BufferedBody) WriteAsJSON(obj any) (err error) {
	b.mutex.Lock()
	b.buffer.Reset()
	err = json.NewEncoder(b.buffer).Encode(obj)
	b.mutex.Unlock()

	return
}

// ReadAsString implements contracts.Body.
func (b *BufferedBody) ReadAsString() (str string, err error) {
	b.mutex.RLock()
	str = b.buffer.String()
	b.mutex.RUnlock()

	return
}

// WriteAsString implements contracts.Body.
func (b *BufferedBody) WriteAsString(body string) (err error) {
	b.mutex.Lock()
	b.buffer.Reset()
	_, err = b.buffer.WriteString(body)
	b.mutex.Unlock()

	return
}

// ReadAsXML implements contracts.Body.
func (b *BufferedBody) ReadAsXML(obj any) (err error) {
	b.mutex.RLock()
	err = xml.Unmarshal(b.buffer.Bytes(), obj)
	b.mutex.RUnlock()

	return
}

// WriteAsXML implements contracts.Body.
func (b *BufferedBody) WriteAsXML(obj any) (err error) {
	b.mutex.Lock()
	b.buffer.Reset()
	err = xml.NewEncoder(b.buffer).Encode(obj)
	b.mutex.Unlock()

	return
}

// Set implements contracts.Body.
func (b *BufferedBody) Set(body io.Reader) (err error) {
	b.mutex.Lock()
	b.buffer.Reset()
	_, err = io.Copy(b.buffer, body)
	b.mutex.Unlock()

	return
}

// Unwrap implements contracts.Body.
func (b *BufferedBody) Unwrap() (body io.Reader) {
	b.mutex.RLock()
	body = bytes.NewReader(b.buffer.Bytes())
	b.mutex.RUnlock()

	return
}

func newBufferedBody() *BufferedBody {
	return &BufferedBody{
		buffer: bytes.NewBuffer(nil),
		mutex:  sync.RWMutex{},
	}
}
