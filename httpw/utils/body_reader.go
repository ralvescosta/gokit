package utils

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-playground/validator/v10"
	"github.com/ralvescosta/gokit/httpw/viewmodels"
)

func BodyReader[B any](reader io.ReadCloser, body *B) *viewmodels.HTTPError {
	err := json.NewDecoder(reader).Decode(body)

	if err != nil {
		return viewmodels.BadRequest("failure to read body")
	}

	val := validator.New()
	err = val.Struct(body)

	if err == nil {
		return nil
	}

	validationErrors := err.(validator.ValidationErrors)
	messages := make(map[string]string, len(validationErrors))

	for _, validationErr := range validationErrors {
		message := ""
		if validationErr.Tag() == "required" {
			message = fmt.Sprintf("%s is required", validationErr.Field())
		} else {
			message = fmt.Sprintf("%s invalid %s", validationErr.Field(), validationErr.Tag())
		}
		messages[validationErr.Field()] = message
	}

	return viewmodels.UnformattedBody(messages)
}
