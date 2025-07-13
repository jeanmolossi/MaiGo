package async

import "errors"

var (
	ErrEmptyRequest      = errors.New("empty request to dispatch")
	ErrNilRequestBuilder = errors.New("nil request builder is not allowed")
	ErrNoRequests        = errors.New("no requests to send")
)
