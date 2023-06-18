package viewmodels

import (
	"encoding/json"
	"net/http"
)

type (

	// HTTPError
	HTTPError struct {
		StatusCode int `json:"status_code"`
		Message    any `json:"message"`
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

func BadRequest(message any) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}

func InternalError(message any) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
	}
}

func NotImplementedYet() *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusNotImplemented,
		Message:    "resource was not implemented yet",
	}
}

func UnformattedBody() *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusBadRequest,
		Message:    "unformatted body",
	}
}

func InvalidBody(message any) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}

func Conflict(message any) *HTTPError {
	return &HTTPError{
		StatusCode: http.StatusConflict,
		Message:    message,
	}
}
