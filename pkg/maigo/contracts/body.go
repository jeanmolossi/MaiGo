package contracts

import "io"

type Body interface {
	io.ReadCloser
	ReadAsJSON(obj any) error
	WriteAsJSON(obj any) error
	ReadAsXML(obj any) error
	WriteAsXML(obj any) error
	ReadAsString() (string, error)
	WriteAsString(body string) error
	Set(body io.Reader) error
	Unwrap() io.Reader
}
