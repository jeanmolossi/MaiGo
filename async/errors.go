package async

import "errors"

var (
	ErrEmptyRequest      = errors.New("empty request to dispatch")
)
