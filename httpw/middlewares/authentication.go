package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ralvescosta/gokit/auth"
	"github.com/ralvescosta/gokit/httpw"
	"github.com/ralvescosta/gokit/httpw/viewmodels"
	"github.com/ralvescosta/gokit/logging"
	"go.uber.org/zap"
)

type (
	Authorization interface {
		Handle(next http.Handler) http.Handler
	}

	authorization struct {
		logger       logging.Logger
		tokenManager auth.IdentityManager
	}
)

func NewAuthorization(logger logging.Logger, tokenManager auth.IdentityManager) Authorization {
	return &authorization{logger, tokenManager}
}

func (a *authorization) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		if authorization == "" {
			msg := "token was not provided"

			a.logger.Error(httpw.Message(msg))

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(viewmodels.HTTPError{
				StatusCode: http.StatusUnauthorized,
				Message:    msg,
			})

			return
		}

		part := strings.Split(authorization, " ")
		if len(part) < 2 || part[0] != "Bearer" || part[1] == "" {
			msg := "unformatted token"

			a.logger.Error(httpw.Message(msg))

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(viewmodels.HTTPError{
				StatusCode: http.StatusUnauthorized,
				Message:    msg,
			})

			return
		}

		session, err := a.tokenManager.Validate(r.Context(), part[1])
		if err != nil {
			a.logger.Error(httpw.Message("failure to validate the token"), zap.Error(err))

			viewmodels.NewResponseBuilder(w).BadRequest().Message(err.Error()).Build()

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(viewmodels.HTTPError{
				StatusCode: http.StatusUnauthorized,
				Message:    err.Error(),
			})

			return
		}

		ctx := context.WithValue(r.Context(), "session", session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
