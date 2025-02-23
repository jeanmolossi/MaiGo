package maigo

import "errors"

var (
	ErrEmptyBaseURL      = errors.New("empty base URL is not allowed")
	ErrParseURL          = errors.New("failed to parse URL")
	ErrClientValidation  = errors.New("invalid client attributes")
	ErrRequestValidation = errors.New("invalid client attributes")
	ErrCreateRequest     = errors.New("failed to create request")
)
