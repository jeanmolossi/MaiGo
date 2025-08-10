package maigo

import (
	"errors"
	"net/url"
	"strings"

	"github.com/jeanmolossi/MaiGo/pkg/maigo/contracts"
)

var _ contracts.BuilderRequestQuery[contracts.RequestBuilder] = (*RequestQueryBuilder)(nil)

type RequestQueryBuilder struct {
	parent *RequestBuilder
	config *RequestConfigBase
}

func (r *RequestBuilder) Query() contracts.BuilderRequestQuery[contracts.RequestBuilder] {
	return &RequestQueryBuilder{
		parent: r,
		config: r.request.config,
	}
}

// AddParam implements contracts.BuilderRequestQuery.
func (r *RequestQueryBuilder) AddParam(key string, value string) contracts.RequestBuilder {
	r.config.searchParams.Add(key, value)
	return r.parent
}

// AddParams implements contracts.BuilderRequestQuery.
func (r *RequestQueryBuilder) AddParams(params contracts.Params) contracts.RequestBuilder {
	for key, value := range params {
		r.AddParam(key, value)
	}

	return r.parent
}

// AddRawString implements contracts.BuilderRequestQuery.
func (r *RequestQueryBuilder) AddRawString(raw string) contracts.RequestBuilder {
	actual := r.config.searchParams.Encode()

	if actual == "" {
		r.SetRawString(raw)
		return r.parent
	}

	var err error

	raw = "&" + trimQueryPrefix(raw)

	r.config.searchParams, err = url.ParseQuery(strings.TrimSpace(actual + raw))
	if err != nil {
		r.config.validations.Add(errors.Join(ErrAddingRawQueryToActualQuery, err))
	}

	return r.parent
}

// SetParam implements contracts.BuilderRequestQuery.
func (r *RequestQueryBuilder) SetParam(key string, value string) contracts.RequestBuilder {
	r.config.searchParams.Set(key, value)
	return r.parent
}

// SetParams implements contracts.BuilderRequestQuery.
func (r *RequestQueryBuilder) SetParams(params contracts.Params) contracts.RequestBuilder {
	for key, value := range params {
		r.SetParam(key, value)
	}

	return r.parent
}

// SetRawString implements contracts.BuilderRequestQuery.
func (r *RequestQueryBuilder) SetRawString(raw string) contracts.RequestBuilder {
	raw = trimQueryPrefix(raw)

	newSearchParams, err := url.ParseQuery(strings.TrimSpace(raw))
	if err != nil {
		r.config.validations.Add(errors.Join(ErrSettingRawQuery, err))
		return r.parent
	}

	r.config.searchParams = newSearchParams

	return r.parent
}

func trimQueryPrefix(q string) string {
	for _, prefix := range []string{"&", "?"} {
		q = strings.TrimPrefix(q, prefix)
	}

	return q
}
