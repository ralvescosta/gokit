package viewmodels

import (
	"encoding/json"
	"net/http"
)

type (
	HTTPErrorMessage interface {
		string | map[string]string
	}

	// HTTPError
	HTTPError[M HTTPErrorMessage] struct {
		StatusCode int `json:"status_code"`
		Message    M   `json:"message"`
	}
)

func (h *HTTPError[HTTPErrorMessage]) ToBuffer() []byte {
	b, _ := json.Marshal(h)

	return b
}

func BadRequest[M HTTPErrorMessage](message M) *HTTPError[M] {
	return &HTTPError[M]{
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}

func InternalError[M HTTPErrorMessage](message M) *HTTPError[M] {
	return &HTTPError[M]{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
	}
}

func NotImplementedYet() *HTTPError[string] {
	return &HTTPError[string]{
		StatusCode: http.StatusNotImplemented,
		Message:    "resource was not implemented yet",
	}
}

func UnformattedBody() *HTTPError[string] {
	return &HTTPError[string]{
		StatusCode: http.StatusBadRequest,
		Message:    "unformatted body",
	}
}

func InvalidBody[M HTTPErrorMessage](message M) *HTTPError[M] {
	return &HTTPError[M]{
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}

func Conflict[M HTTPErrorMessage](message M) *HTTPError[M] {
	return &HTTPError[M]{
		StatusCode: http.StatusConflict,
		Message:    message,
	}
}
