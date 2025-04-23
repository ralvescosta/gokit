// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package httpw provides HTTP wrapper utilities for building robust HTTP services.
// It offers a set of tools to easily create and manage HTTP servers, define routes,
// handle middleware, validate requests, and format responses.
package httpw

// Message formats a message with the httpw package prefix.
// This function is useful for consistent logging across the httpw package.
func Message(msg string) string {
	return "[gokit:httpw] " + msg
}
