package middleware

import (
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"

	"github.com/f-code-club/rode-battle-api/internal/shared"
	"github.com/f-code-club/rode-battle-api/internal/shared/headers"
)

type ContextKey string

const (
	bearerPrefix            = "Bearer "
	AccountIDKey ContextKey = "id"
)

func NewParseToken(api huma.API, accessTokenSvc *shared.TokenService) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		authHeader := ctx.Header(headers.Authorization)
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "missing or invalid Authorization header")
			return
		}

		tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
		userId, err := accessTokenSvc.ParseToken(tokenStr)
		if err != nil {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "invalid token", err)
			return
		}

		ctx = huma.WithValue(ctx, AccountIDKey, userId)
		next(ctx)
	}
}
