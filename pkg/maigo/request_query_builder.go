package maigo

import (
	"errors"
	"net/url"
	"strings"

	"github.com/jeanmolossi/maigo/pkg/maigo/contracts"
)

var _ contracts.BuilderRequestQuery[contracts.RequestBuilder] = (*RequestQueryBuilder)(nil)

const maxQueryLen = 8 << 10 // 8KiB

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
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return r.parent
	}

	raw = trimQueryPrefix(raw)

	if len(raw) > maxQueryLen || hasCTL(raw) {
		r.config.validations.Add(ErrInvalidQueryString)
		return r.parent
	}

	inc, err := url.ParseQuery(raw)
	if err != nil {
		r.config.validations.Add(errors.Join(ErrAddingRawQueryToActualQuery, err))
		return r.parent
	}

	if r.config.searchParams == nil {
		r.config.searchParams = make(url.Values, len(inc))
	}

	for k, vs := range inc {
		for _, v := range vs {
			r.config.searchParams.Add(k, v)
		}
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
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return r.parent
	}

	raw = trimQueryPrefix(raw)

	if len(raw) > maxQueryLen || hasCTL(raw) {
		r.config.validations.Add(ErrInvalidQueryString)
		return r.parent
	}

	newSearchParams, err := url.ParseQuery(raw)
	if err != nil {
		r.config.validations.Add(errors.Join(ErrSettingRawQuery, err))
		return r.parent
	}

	r.config.searchParams = newSearchParams

	return r.parent
}

func trimQueryPrefix(q string) string {
	for len(q) > 0 && (q[0] == '?' || q[0] == '&') {
		q = q[1:]
	}

	return q
}

func hasCTL(s string) bool {
	for _, r := range s {
		// RFC3986 CTLs out of percent-encoding should be avoided
		// (0x00-0x1F and 0x7F)
		if r < 0x20 || r == 0x7f {
			return true
		}
	}

	return false
}
