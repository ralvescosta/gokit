// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/ralvescosta/gokit/logging"

	"github.com/ralvescosta/gokit/httpw/viewmodels"
)

type (
	// BodyValidator defines the interface for validating request bodies.
	BodyValidator interface {
		Validate(body any) *viewmodels.HTTPError
	}

	// bodyValidator implements the BodyValidator interface.
	bodyValidator struct {
		logger logging.Logger
	}
)

func NewBodyValidator(logging logging.Logger) BodyValidator {
	return &bodyValidator{logging}
}

func (b *bodyValidator) Validate(body any) *viewmodels.HTTPError {
	val := validator.New()
	err := val.Struct(body)

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

	res := viewmodels.BadRequest(messages)
	res.Message = "invalid body"

	return res
}
