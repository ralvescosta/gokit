// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package httpw provides HTTP wrapper utilities for building robust HTTP services.
package httpw

import "errors"

// Common HTTP error variables used throughout the httpw package
var (
	// ErrInvalidHTTPMethod indicates that the provided HTTP method is not supported.
	// This error is returned when attempting to register a route with an unsupported HTTP method.
	ErrInvalidHTTPMethod = errors.New("invalid http method")

	// ErrHTTPMethodMethodIsRequired indicates that an HTTP method is required but not provided.
	// This error is returned when attempting to register a route without specifying the HTTP method.
	ErrHTTPMethodMethodIsRequired = errors.New("http method is required")
)
