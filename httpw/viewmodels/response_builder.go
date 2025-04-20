// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package viewmodels

import (
	"encoding/json"
	"net/http"
)

type ResponseBuilder struct {
	writer       http.ResponseWriter
	statusCode   int
	body         any
	errorDetails any
	errorMessage string
	header       map[string]string
}

func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{header: map[string]string{}}
}

func (b *ResponseBuilder) Writer(writer http.ResponseWriter) *ResponseBuilder {
	b.writer = writer
	return b
}

func (b *ResponseBuilder) Ok() *ResponseBuilder {
	b.statusCode = http.StatusOK
	return b
}

func (b *ResponseBuilder) Created() *ResponseBuilder {
	b.statusCode = http.StatusCreated
	return b
}

func (b *ResponseBuilder) NoContent() *ResponseBuilder {
	b.statusCode = http.StatusNoContent
	return b
}

func (b *ResponseBuilder) BadRequest() *ResponseBuilder {
	b.statusCode = http.StatusBadRequest
	return b
}

func (b *ResponseBuilder) Unauthorized() *ResponseBuilder {
	b.statusCode = http.StatusUnauthorized
	return b
}

func (b *ResponseBuilder) Forbidden() *ResponseBuilder {
	b.statusCode = http.StatusForbidden
	return b
}

func (b *ResponseBuilder) NotFound() *ResponseBuilder {
	b.statusCode = http.StatusNotFound
	return b
}

func (b *ResponseBuilder) Conflict() *ResponseBuilder {
	b.statusCode = http.StatusConflict
	return b
}

func (b *ResponseBuilder) InternalError() *ResponseBuilder {
	b.statusCode = http.StatusInternalServerError
	return b
}

func (b *ResponseBuilder) Message(m string) *ResponseBuilder {
	b.errorMessage = m
	return b
}

func (b *ResponseBuilder) Details(m any) *ResponseBuilder {
	b.errorDetails = m
	return b
}

func (b *ResponseBuilder) JSON(body any) *ResponseBuilder {
	b.body = body
	return b
}

func (b *ResponseBuilder) Header(key, value string) *ResponseBuilder {
	b.header[key] = value
	return b
}

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
			Details:    b.Details,
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
