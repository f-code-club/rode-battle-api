package middleware

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/f-code-club/rode-battle-api/internal/auth/service"
	"github.com/google/uuid"
)

func NewRequireRole(api huma.API, authSvc service.Service, roles ...service.Role) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		accountID, ok := ctx.Context().Value(AccountIDKey).(uuid.UUID)
		if !ok {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "unauthorized: account not found in context")
			return
		}

		allowed, err := authSvc.HasRole(ctx.Context(), accountID, roles...)
		if err != nil {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "unauthorized: failed to verify role", err)
			return
		}

		if !allowed {
			_ = huma.WriteErr(api, ctx, http.StatusForbidden, "forbidden: insufficient permissions")
			return
		}

		next(ctx)
	}
}
