package contracts

// Params represents a set of query parameters to be added to a request URL.
type Params map[string]string

// BuilderRequestQuery builds the query string for a request. It supports
// adding and setting parameters individually or in bulk and also accepts raw
// query strings.
//
// Example:
//
//	req.Query().
//	        AddParam("page", "2").
//	        AddParams(Params{"sort": "name"})
type BuilderRequestQuery[T any] interface {
	// AddParam appends a single query parameter.
	AddParam(key, value string) T
	// AddParams appends multiple query parameters.
	AddParams(params Params) T
	// SetParam replaces the value for the given key.
	SetParam(key, value string) T
	// SetParams replaces the entire set of parameters.
	SetParams(params Params) T
	// AddRawString appends raw query string data.
	AddRawString(raw string) T
	// SetRawString replaces the raw query string.
	SetRawString(raw string) T
}
