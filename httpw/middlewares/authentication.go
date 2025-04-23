// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package middlewares provides HTTP middleware components for the httpw package.
// Middlewares can be used to add cross-cutting concerns like authentication,
// logging, or request tracing to HTTP request handling.
package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ralvescosta/gokit/auth"
	"github.com/ralvescosta/gokit/logging"

	"github.com/ralvescosta/gokit/httpw"
	"github.com/ralvescosta/gokit/httpw/viewmodels"
)

type (
	// Authorization defines the interface for handling authorization middleware.
	// It provides functionality to authenticate incoming HTTP requests based on JWT tokens.
	Authorization interface {
		// Handle returns an http.Handler middleware that performs token validation.
		// It extracts the token from the Authorization header, validates it, and adds
		// the authenticated session to the request context.
		Handle(next http.Handler) http.Handler
	}

	// authorization implements the Authorization interface.
	authorization struct {
		logger       logging.Logger
		tokenManager auth.IdentityManager
	}
)

// NewAuthorization creates a new Authorization middleware with the provided logger and token manager.
func NewAuthorization(logger logging.Logger, tokenManager auth.IdentityManager) Authorization {
	return &authorization{logger, tokenManager}
}

// Handle returns an http.Handler middleware that performs token validation.
// It extracts the token from the Authorization header, validates it, and adds
// the authenticated session to the request context.
func (a *authorization) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		if authorization == "" {
			msg := "token was not provided"

			a.logger.Error(httpw.Message(msg))

			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(viewmodels.HTTPError{
				StatusCode: http.StatusUnauthorized,
				Message:    msg,
			})

			return
		}

		part := strings.Split(authorization, " ")
		if len(part) < 2 || part[0] != "Bearer" || part[1] == "" {
			a.handleError(w, "unformatted token", http.StatusUnauthorized)
			return
		}

		session, err := a.tokenManager.Validate(r.Context(), part[1])
		if err != nil {
			a.handleError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), &auth.Claims{}, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// handleError logs the error and writes an HTTP error response.
func (a *authorization) handleError(w http.ResponseWriter, msg string, statusCode int) {
	a.logger.Error(httpw.Message(msg))
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(viewmodels.HTTPError{
		StatusCode: statusCode,
		Message:    msg,
	})
}
