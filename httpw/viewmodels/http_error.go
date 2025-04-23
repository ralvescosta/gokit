// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package viewmodels provides standardized response structures for HTTP handlers.
// It offers consistent error representations and response building utilities.
package viewmodels

import (
	"encoding/json"
	"net/http"
)

type (
	// HTTPError represents an HTTP error response with a status code, message, and details.
	// It provides a standardized structure for error responses across the API.
	HTTPError struct {
		StatusCode int    `json:"status_code" example:"400"`         // HTTP status code
		Message    string `json:"message"     example:"bad request"` // Error message
		Details    any    `json:"details"`                           // Additional error details
	}
)

// ToBuffer serializes the HTTPError to a JSON byte array.
func (h *HTTPError) ToBuffer() []byte {
	b, _ := json.Marshal(h)
	return b
}

// ToString serializes the HTTPError to a JSON string.
func (h *HTTPError) ToString() string {
	b, _ := json.Marshal(h)
	return string(b)
}

// BadRequest creates a new HTTPError with 400 Bad Request status code.
func BadRequest(details any) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusBadRequest,
		Message:    "Bad Request",
		Details:    details,
	}
}

// InternalError creates a new HTTPError with 500 Internal Server Error status code.
func InternalError(details any) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Error",
		Details:    details,
	}
}

// NotImplementedYet creates a new HTTPError with 501 Not Implemented status code.
func NotImplementedYet() *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusNotImplemented,
		Message:    "resource was not implemented yet",
	}
}

// UnformattedBody creates a new HTTPError for malformed request body with 400 Bad Request status code.
func UnformattedBody(details any) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusBadRequest,
		Message:    "unformatted body",
		Details:    details,
	}
}

// Conflict creates a new HTTPError with 409 Conflict status code.
func Conflict(details any) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusConflict,
		Details:    details,
	}
}
