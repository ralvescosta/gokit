package httpw

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func Validate[T any](reader io.ReadCloser, body *T) *HTTPError {
	err := json.NewDecoder(reader).Decode(body)
	if err != nil {
		return &HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    "unformatted body",
			Details:    "wrong json format",
		}
	}

	validate := validator.New()
	err = validate.Struct(body)
	if err != nil {
		details := map[string]map[string]string{}

		for _, err := range err.(validator.ValidationErrors) {
			details[err.Field()] = map[string]string{err.Tag(): err.Param()}
		}

		return &HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    "unformatted body",
			Details:    details,
		}
	}

	return nil
}
