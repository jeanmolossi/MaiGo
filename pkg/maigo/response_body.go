package maigo

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
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

// AsJSON decodes the body as JSON into the provided value.
//
// Deprecated: Use [AsJSON] generic function instead.
func (r *ResponseBody) AsJSON(v any) error {
	raw, err := AsJSON[json.RawMessage](r)
	if err != nil {
		return err
	}

	return json.Unmarshal(raw, v)
}

// AsJSON decodes the body as JSON into a new instance of T.
func AsJSON[T any](r *ResponseBody) (T, error) {
	defer r.Close()

	var v T
	if err := r.body.ReadAsJSON(&v); err != nil {
		var zero T
		return zero, fmt.Errorf("failed reading as JSON: %w", err)
	}

	return v, nil
}

// AsString implements contracts.ResponseFluentBody.
func (r *ResponseBody) AsString() (string, error) {
	defer r.Close()
	return r.body.ReadAsString()
}

// AsXML implements contracts.ResponseFluentBody.
//
// Deprecated: Use [AsXML] generic function instead.
func (r *ResponseBody) AsXML(v any) error {
	raw, err := AsXML[[]byte](r)
	if err != nil {
		return err
	}

	return xml.Unmarshal(raw, v)
}

// AsXML decodes the body as XML into a new instance of T.
func AsXML[T any](r *ResponseBody) (T, error) {
	defer r.Close()

	var v T
	if err := r.body.ReadAsXML(&v); err != nil {
		var zero T
		return zero, fmt.Errorf("failed reading as XML: %w", err)
	}

	return v, nil
}

// Close implements contracts.ResponseFluentBody.
func (r *ResponseBody) Close() {
	_ = r.body.Close()
}

// Raw implements contracts.ResponseFluentBody.
func (r *ResponseBody) Raw() io.Closer {
	return r.body
}
