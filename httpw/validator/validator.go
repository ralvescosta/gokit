// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package validator provides request validation utilities for HTTP requests.
// It offers tools to validate incoming request bodies against structural and semantic rules.
package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/ralvescosta/gokit/logging"

	"github.com/ralvescosta/gokit/httpw/viewmodels"
)

type (
	// BodyValidator defines the interface for validating request bodies.
	// It provides functionality to validate structured data in HTTP request bodies
	// against validation rules defined using struct tags.
	BodyValidator interface {
		// Validate checks if the provided body adheres to validation rules.
		// It returns an HTTPError if validation fails, or nil if the body is valid.
		Validate(body any) *viewmodels.HTTPError
	}

	// bodyValidator implements the BodyValidator interface.
	bodyValidator struct {
		logger logging.Logger
	}
)

// NewBodyValidator creates a new BodyValidator with the provided logger.
func NewBodyValidator(logging logging.Logger) BodyValidator {
	return &bodyValidator{logging}
}

// Validate checks if the provided body adheres to validation rules defined by struct tags.
// It returns an HTTPError with detailed validation error messages if validation fails,
// or nil if the body is valid.
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
