package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-fuego/fuego"
)

//nolint:unused // will be used by auth middleware context
type contextKey string

const (
	//nolint:unused // will be used by auth middleware context
	authHeaderKey = "Authorization"
	//nolint:unused // will be used by auth middleware context
	bearerPrefix = "Bearer "
	//nolint:unused // will be used by auth middleware context
	userKey contextKey = "user"
)

//nolint:unused // will be used by auth middleware context
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
		userId, err := s.auth.TokenService.ParseToken(tokenStr)
		if err != nil {
			fuego.SendJSONError(w, nil, fuego.UnauthorizedError{
				Detail: "invalid token",
				Err:    err,
			})
			return
		}
		ctx := context.WithValue(r.Context(), userKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
