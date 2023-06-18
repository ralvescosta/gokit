package viewmodels

import (
	"encoding/json"
	"net/http"
)

type (

	// HTTPError
	HTTPError struct {
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
		Details    any    `json:"details"`
	}
)

// func AnyToErrorMessage(message any) M {
// 	if m, ok := message.(string); ok {
// 		return m
// 	}

// 	m, _ := message.(mapstring)
// 	return m
// }

func (h *HTTPError) ToBuffer() []byte {
	b, _ := json.Marshal(h)

	return b
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
