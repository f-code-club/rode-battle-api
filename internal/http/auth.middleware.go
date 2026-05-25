package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-fuego/fuego"
)

type contextKey string

const (
	authHeaderKey = "Authorization"
	bearerPrefix  = "Bearer "

	userKey contextKey = "user"
)

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(authHeaderKey)
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			fuego.SendJSONError(w, nil, fuego.UnauthorizedError{
				Detail: "missing or invalid Authorization header",
			})
			return
		}

		tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
		userId, err := s.auth.ParseToken(tokenStr)
		if err != nil {
			fuego.SendJSONError(w, nil, fuego.UnauthorizedError{
				Detail: "invalid token",
			})
			return
		}
		ctx := context.WithValue(r.Context(), userKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
