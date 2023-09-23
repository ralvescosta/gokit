package viewmodels

import (
	"encoding/json"
	"net/http"
)

type responseBuilder struct {
	writer       http.ResponseWriter
	statusCode   int
	body         any
	errorDetails any
	errorMessage string
	header       map[string]string
}

func NewResponseBuilder() *responseBuilder {
	return &responseBuilder{header: map[string]string{}}
}

func (b *responseBuilder) Writer(writer http.ResponseWriter) *responseBuilder {
	b.writer = writer
	return b
}

func (b *responseBuilder) Ok() *responseBuilder {
	b.statusCode = http.StatusOK
	return b
}

func (b *responseBuilder) Created() *responseBuilder {
	b.statusCode = http.StatusCreated
	return b
}

func (b *responseBuilder) NoContent() *responseBuilder {
	b.statusCode = http.StatusNoContent
	return b
}

func (b *responseBuilder) BadRequest() *responseBuilder {
	b.statusCode = http.StatusBadRequest
	return b
}

func (b *responseBuilder) Unauthorized() *responseBuilder {
	b.statusCode = http.StatusUnauthorized
	return b
}

func (b *responseBuilder) Forbidden() *responseBuilder {
	b.statusCode = http.StatusForbidden
	return b
}

func (b *responseBuilder) NotFound() *responseBuilder {
	b.statusCode = http.StatusNotFound
	return b
}

func (b *responseBuilder) Conflict() *responseBuilder {
	b.statusCode = http.StatusConflict
	return b
}

func (b *responseBuilder) InternalError() *responseBuilder {
	b.statusCode = http.StatusInternalServerError
	return b
}

func (b *responseBuilder) Message(m string) *responseBuilder {
	b.errorMessage = m
	return b
}

func (b *responseBuilder) Details(m any) *responseBuilder {
	b.errorDetails = m
	return b
}

func (b *responseBuilder) JSON(body any) *responseBuilder {
	b.body = body
	return b
}

func (b *responseBuilder) Header(key, value string) *responseBuilder {
	b.header[key] = value
	return b
}

func (b *responseBuilder) Build() {
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

		b.writer.Write(err.ToBuffer())
		return
	}

	bytes, _ := json.Marshal(b.body)
	b.writer.Write(bytes)
}
