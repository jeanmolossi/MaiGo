package client

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"sync"

	"github.com/jeanmolossi/MaiGo/pkg/client/contracts"
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
		mytex  sync.RWMutex
	}
)

// Close implements contracts.Body.
func (u *UnbufferedBody) Close() error {
	panic("unimplemented")
}

// Read implements contracts.Body.
func (u *UnbufferedBody) Read(p []byte) (n int, err error) {
	panic("unimplemented")
}

// ReadAsJSON implements contracts.Body.
func (u *UnbufferedBody) ReadAsJSON(obj any) error {
	panic("unimplemented")
}

// ReadAsString implements contracts.Body.
func (u *UnbufferedBody) ReadAsString() (string, error) {
	panic("unimplemented")
}

// ReadAsXML implements contracts.Body.
func (u *UnbufferedBody) ReadAsXML(obj any) error {
	panic("unimplemented")
}

// Set implements contracts.Body.
func (u *UnbufferedBody) Set(body io.Reader) error {
	panic("unimplemented")
}

// Unwrap implements contracts.Body.
func (u *UnbufferedBody) Unwrap() io.Reader {
	panic("unimplemented")
}

// WriteAsJSON implements contracts.Body.
func (u *UnbufferedBody) WriteAsJSON(obj any) error {
	panic("unimplemented")
}

// WriteAsString implements contracts.Body.
func (u *UnbufferedBody) WriteAsString(body string) error {
	panic("unimplemented")
}

// WriteAsXML implements contracts.Body.
func (u *UnbufferedBody) WriteAsXML(obj any) error {
	panic("unimplemented")
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
	err = json.NewDecoder(b.buffer).Decode(obj)
	b.mutex.RUnlock()

	return
}

// WriteAsJSON implements contracts.Body.
func (b *BufferedBody) WriteAsJSON(obj any) (err error) {
	b.mutex.Lock()
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
	_, err = b.buffer.WriteString(body)
	b.mutex.Unlock()

	return
}

// ReadAsXML implements contracts.Body.
func (b *BufferedBody) ReadAsXML(obj any) (err error) {
	b.mutex.RLock()
	err = xml.NewDecoder(b.buffer).Decode(obj)
	b.mutex.RUnlock()

	return
}

// WriteAsXML implements contracts.Body.
func (b *BufferedBody) WriteAsXML(obj any) (err error) {
	b.mutex.Lock()
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
