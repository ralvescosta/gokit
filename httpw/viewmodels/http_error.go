// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package viewmodels

import (
	"encoding/json"
	"net/http"
)

type (

	// HTTPError
	HTTPError struct {
		StatusCode int    `json:"status_code" example:"400"`
		Message    string `json:"message"     example:"bad request"`
		Details    any    `json:"details"`
	}
)

func (h *HTTPError) ToBuffer() []byte {
	b, _ := json.Marshal(h)

	return b
}

func (h *HTTPError) ToString() string {
	b, _ := json.Marshal(h)

	return string(b)
}

func BadRequest(details any) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusBadRequest,
		Message:    "Bad Request",
		Details:    details,
	}
}

func InternalError(details any) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Error",
		Details:    details,
	}
}

func NotImplementedYet() *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusNotImplemented,
		Message:    "resource was not implemented yet",
	}
}

func UnformattedBody(details any) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusBadRequest,
		Message:    "unformatted body",
		Details:    details,
	}
}

func Conflict(details any) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusConflict,
		Details:    details,
	}
}
