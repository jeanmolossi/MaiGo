package maigo

import (
	"errors"
	"io"

	"github.com/jeanmolossi/MaiGo/pkg/maigo/contracts"
)

var _ contracts.BuilderRequestBody[contracts.RequestBuilder] = (*RequestBodyBuilder)(nil)

type RequestBodyBuilder struct {
	parent *RequestBuilder
	config *RequestConfigBase
}

func (r *RequestBuilder) Body() contracts.BuilderRequestBody[contracts.RequestBuilder] {
	return &RequestBodyBuilder{
		parent: r,
		config: r.request.config,
	}
}

// AsReader implements contracts.BuilderRequestBody.
func (r *RequestBodyBuilder) AsReader(body io.Reader) contracts.RequestBuilder {
	err := r.config.body.Set(body)
	if err != nil {
		r.config.validations.Add(errors.Join(ErrToSetBody, err))
	}

	return r.parent
}

// AsJSON implements contracts.BuilderRequestBody.
func (r *RequestBodyBuilder) AsJSON(obj any) contracts.RequestBuilder {
	err := r.config.body.WriteAsJSON(obj)
	if err != nil {
		r.config.validations.Add(errors.Join(ErrToMarshalJSON, err))
	}

	return r.parent
}

// AsString implements contracts.BuilderRequestBody.
func (r *RequestBodyBuilder) AsString(body string) contracts.RequestBuilder {
	err := r.config.body.WriteAsString(body)
	if err != nil {
		r.config.validations.Add(errors.Join(ErrToSetBody, err))
	}

	return r.parent
}

// AsXML implements contracts.BuilderRequestBody.
func (r *RequestBodyBuilder) AsXML(obj any) contracts.RequestBuilder {
	err := r.config.body.WriteAsXML(obj)
	if err != nil {
		r.config.validations.Add(errors.Join(ErrToMarshalXML, err))
	}

	return r.parent
}
