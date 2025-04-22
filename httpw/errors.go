// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package httpw

import "errors"

var (
	// ErrInvalidHTTPMethod indicates that the provided HTTP method is not supported.
	ErrInvalidHTTPMethod = errors.New("invalid http method")

	// ErrHTTPMethodMethodIsRequired indicates that an HTTP method is required but not provided.
	ErrHTTPMethodMethodIsRequired = errors.New("http method is required")
)
