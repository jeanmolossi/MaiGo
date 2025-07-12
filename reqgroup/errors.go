package reqgroup

import "errors"

var (
	ErrNilRequestBuilder = errors.New("nil request builder is not allowed")
	ErrNoRequests        = errors.New("no requests to send")
)
