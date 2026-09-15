package middleware

import (
	"net/http"

	"github.com/f-code-club/rode-battle-api/internal/auth/service"
	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
)

func NewRequireRole(authSvc service.Service, roles ...service.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accountID, ok := r.Context().Value(AccountIDKey).(uuid.UUID)
			if !ok {
				fuego.SendJSONError(w, nil, fuego.UnauthorizedError{
					Detail: "unauthorized: account not found in context",
				})
				return
			}

			allowed, err := authSvc.HasRole(r.Context(), accountID, roles...)
			if err != nil {
				fuego.SendJSONError(w, nil, fuego.UnauthorizedError{
					Detail: "unauthorized: failed to verify role",
					Err:    err,
				})
				return
			}

			if !allowed {
				fuego.SendJSONError(w, nil, fuego.ForbiddenError{
					Detail: "forbidden: insufficient permissions",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
