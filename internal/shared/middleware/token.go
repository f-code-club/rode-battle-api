package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-fuego/fuego"

	"github.com/f-code-club/rode-battle-api/internal/shared"
	"github.com/f-code-club/rode-battle-api/internal/shared/headers"
)

const (
	bearerPrefix = "Bearer "
	AccountIDKey = "id"
)

type ParseTokenMiddlewareBuilder struct {
	Service *shared.TokenService
}

func (b ParseTokenMiddlewareBuilder) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(headers.Authorization)
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			fuego.SendJSONError(w, nil, fuego.UnauthorizedError{
				Detail: "missing or invalid Authorization header",
			})
			return
		}

		tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
		userId, err := b.Service.ParseToken(tokenStr)
		if err != nil {
			fuego.SendJSONError(w, nil, fuego.UnauthorizedError{
				Detail: "invalid token",
				Err:    err,
			})
			return
		}
		ctx := context.WithValue(r.Context(), AccountIDKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
