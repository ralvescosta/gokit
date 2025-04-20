// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package httpw

import "errors"

var (
	ErrInvalidHTTPMethod          = errors.New("invalid http method")
	ErrHTTPMethodMethodIsRequired = errors.New("http method is required")
)
