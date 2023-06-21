package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/ralvescosta/gokit/auth"
	"github.com/ralvescosta/gokit/httpw"
	"github.com/ralvescosta/gokit/httpw/viewmodels"
	"github.com/ralvescosta/gokit/logging"
	"go.uber.org/zap"
)

type (
	Authorization interface{}

	authorization struct {
		logger       logging.Logger
		tokenManager auth.IdentityManager
	}
)

func NewAuthorization(logger logging.Logger, tokenManager auth.IdentityManager) Authorization {
	return &authorization{logger, tokenManager}
}

func (a *authorization) Mid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			msg := "token was not provided"
			a.logger.Error(httpw.Message(msg))
			viewmodels.NewResponseBuilder(w).BadRequest().Message(msg).Build()
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		part := strings.Split(authorization, " ")
		if len(part) < 2 || part[0] != "Bearer" || part[1] == "" {
			msg := "unformatted token"
			a.logger.Error(httpw.Message(msg))
			viewmodels.NewResponseBuilder(w).BadRequest().Message(msg).Build()
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		session, err := a.tokenManager.Validate(r.Context(), part[1])
		if err != nil {
			a.logger.Error(httpw.Message("failure to validate the token"), zap.Error(err))
			viewmodels.NewResponseBuilder(w).BadRequest().Message(err.Error()).Build()
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "session", session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
