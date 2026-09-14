package middleware

import (
	"context"
	"net/http"

	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
)

type CheckRoleFunc func(ctx context.Context, accountID uuid.UUID) (bool, error)

func NewRequireRole(check CheckRoleFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accountID, ok := r.Context().Value(AccountIDKey).(uuid.UUID)
			if !ok {
				fuego.SendJSONError(w, nil, fuego.UnauthorizedError{
					Detail: "unauthorized: account not found in context",
				})
				return
			}

			allowed, err := check(r.Context(), accountID)
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
