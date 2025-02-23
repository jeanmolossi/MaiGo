package client

import "errors"

var (
	ErrEmptyBaseURL = errors.New("empty base URL is not allowed")
	ErrParseURL     = errors.New("failed to parse URL")
)
