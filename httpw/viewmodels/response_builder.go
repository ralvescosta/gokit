// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package viewmodels

import (
	"encoding/json"
	"net/http"
)

// ResponseBuilder helps in constructing HTTP responses with various configurations.
type ResponseBuilder struct {
	writer       http.ResponseWriter
	statusCode   int
	body         any
	errorDetails any
	errorMessage string
	header       map[string]string
}

// NewResponseBuilder creates a new instance of ResponseBuilder.
func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{header: map[string]string{}}
}

// Writer sets the HTTP response writer.
func (b *ResponseBuilder) Writer(writer http.ResponseWriter) *ResponseBuilder {
	b.writer = writer
	return b
}

// Ok sets the status code to 200 OK.
func (b *ResponseBuilder) Ok() *ResponseBuilder {
	b.statusCode = http.StatusOK
	return b
}

// Created sets the status code to 201 Created.
func (b *ResponseBuilder) Created() *ResponseBuilder {
	b.statusCode = http.StatusCreated
	return b
}

// NoContent sets the status code to 204 No Content.
func (b *ResponseBuilder) NoContent() *ResponseBuilder {
	b.statusCode = http.StatusNoContent
	return b
}

// BadRequest sets the status code to 400 Bad Request.
func (b *ResponseBuilder) BadRequest() *ResponseBuilder {
	b.statusCode = http.StatusBadRequest
	return b
}

// Unauthorized sets the status code to 401 Unauthorized.
func (b *ResponseBuilder) Unauthorized() *ResponseBuilder {
	b.statusCode = http.StatusUnauthorized
	return b
}

// Forbidden sets the status code to 403 Forbidden.
func (b *ResponseBuilder) Forbidden() *ResponseBuilder {
	b.statusCode = http.StatusForbidden
	return b
}

// NotFound sets the status code to 404 Not Found.
func (b *ResponseBuilder) NotFound() *ResponseBuilder {
	b.statusCode = http.StatusNotFound
	return b
}

// Conflict sets the status code to 409 Conflict.
func (b *ResponseBuilder) Conflict() *ResponseBuilder {
	b.statusCode = http.StatusConflict
	return b
}

// InternalError sets the status code to 500 Internal Server Error.
func (b *ResponseBuilder) InternalError() *ResponseBuilder {
	b.statusCode = http.StatusInternalServerError
	return b
}

// Message sets the error message for the response.
func (b *ResponseBuilder) Message(m string) *ResponseBuilder {
	b.errorMessage = m
	return b
}

// Details sets the error details for the response.
func (b *ResponseBuilder) Details(m any) *ResponseBuilder {
	b.errorDetails = m
	return b
}

// JSON sets the response body as JSON.
func (b *ResponseBuilder) JSON(body any) *ResponseBuilder {
	b.body = body
	return b
}

// Header adds a header to the response.
func (b *ResponseBuilder) Header(key, value string) *ResponseBuilder {
	b.header[key] = value
	return b
}

// Build constructs the HTTP response.
func (b *ResponseBuilder) Build() {
	b.writer.WriteHeader(b.statusCode)

	header := b.writer.Header()
	for k, v := range b.header {
		header.Add(k, v)
	}

	if b.statusCode >= 400 {
		err := HTTPError{
			StatusCode: b.statusCode,
			Message:    b.errorMessage,
			Details:    b.errorDetails,
		}

		_, _ = b.writer.Write(err.ToBuffer())
		return
	}

	bytes, err := json.Marshal(b.body)
	if err != nil {
		println(err)
	}

	_, _ = b.writer.Write(bytes)
}
