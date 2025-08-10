package contracts

import "io"

// Body wraps an io.ReadCloser and provides helpers for serializing and
// deserializing HTTP bodies in different formats. Implementations like
// BufferedBody and UnbufferedBody manage the underlying reader and are used
// by both requests and responses.
//
// Example of writing and reading JSON:
//
//	resp, _ := client.POST("/users").
//	        Body().AsJSON(user).
//	        Send()
//	var out User
//	_ = resp.Body().AsJSON(&out)
//
// All implementations are safe for concurrent use.
type Body interface {
	io.ReadCloser
	// ReadAsJSON decodes the body into obj using JSON.
	ReadAsJSON(obj any) error
	// WriteAsJSON encodes obj as JSON and replaces the body contents.
	WriteAsJSON(obj any) error
	// ReadAsXML decodes the body into obj using XML.
	ReadAsXML(obj any) error
	// WriteAsXML encodes obj as XML and replaces the body contents.
	WriteAsXML(obj any) error
	// ReadAsString returns the entire body as a string.
	ReadAsString() (string, error)
	// WriteAsString writes a string to the body.
	WriteAsString(body string) error
	// Set replaces the current body reader with the provided one.
	Set(body io.Reader) error
	// Unwrap exposes the underlying reader without closing it.
	Unwrap() io.Reader
}
