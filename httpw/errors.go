package httpw

import "errors"

var (
	ErrInvalidHTTPMethod          = errors.New("invalid http method")
	ErrHTTPMethodMethodIsRequired = errors.New("http method is required")
)
