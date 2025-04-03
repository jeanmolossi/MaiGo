package maigo

import "errors"

var (
	ErrEmptyBaseURL      = errors.New("empty base URL is not allowed")
	ErrParseURL          = errors.New("failed to parse URL")
	ErrClientValidation  = errors.New("invalid client attributes")
	ErrRequestValidation = errors.New("invalid client attributes")
	ErrCreateRequest     = errors.New("failed to create request")
	ErrParseProxyURL     = errors.New("failed to parse proxy url")
	ErrToSetBody         = errors.New("failed to set body")
	ErrToMarshalJSON     = errors.New("failed to marshal json")
	ErrToMarshalXML      = errors.New("failed to marshal xml")
)
