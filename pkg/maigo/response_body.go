package maigo

import (
	"bytes"
	"fmt"
	"io"

	"github.com/jeanmolossi/MaiGo/pkg/maigo/contracts"
)

var _ contracts.ResponseFluentBody = (*ResponseBody)(nil)

type ResponseBody struct {
	body contracts.Body
}

// AsBytes implements contracts.ResponseFluentBody.
func (r *ResponseBody) AsBytes() ([]byte, error) {
	defer r.Close()

	buf := bytes.NewBuffer(nil)

	_, err := buf.ReadFrom(r.body)
	if err != nil {
		return nil, fmt.Errorf("failed reading as bytes: %w", err)
	}

	return buf.Bytes(), nil
}

// AsJSON implements contracts.ResponseFluentBody.
func (r *ResponseBody) AsJSON(v any) error {
	defer r.Close()
	return r.body.ReadAsJSON(v)
}

// AsString implements contracts.ResponseFluentBody.
func (r *ResponseBody) AsString() (string, error) {
	defer r.Close()
	return r.body.ReadAsString()
}

// AsXML implements contracts.ResponseFluentBody.
func (r *ResponseBody) AsXML(v any) error {
	defer r.Close()
	return r.body.ReadAsXML(v)
}

// Close implements contracts.ResponseFluentBody.
func (r *ResponseBody) Close() {
	_ = r.body.Close()
}

// Raw implements contracts.ResponseFluentBody.
func (r *ResponseBody) Raw() io.Closer {
	return r.body
}
