package middlewares

import (
	"context"
	"net/http"
	"strings"
)

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			w.Write([]byte("token was not provided"))
			http.Error(w, http.StatusText(401), 401)
			return
		}

		part := strings.Split(authorization, " ")
		if len(part) < 2 || part[0] != "Bearer" || part[1] == "" {
			w.Write([]byte("unformatted token"))
			http.Error(w, http.StatusText(401), 401)
			return
		}

		session := map[string]string{}

		ctx := context.WithValue(r.Context(), "session", session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
