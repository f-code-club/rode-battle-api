package http

import (
	"net/http"

	"github.com/f-code-club/rode-battle-api/internal/database"
	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
)

//nolint:unused // will be used by endpoints that require specific roles
func (s *Server) requireRole(roles ...database.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value(userKey).(uuid.UUID)
			if !ok {
				fuego.SendJSONError(w, nil, fuego.UnauthorizedError{
					Detail: "unauthorized: user not found in context",
				})
				return
			}

			allowed, err := s.auth.HasRole(r.Context(), userID.String(), roles...)
			if err != nil {
				fuego.SendJSONError(w, nil, fuego.UnauthorizedError{
					Detail: "unauthorized: failed to verify roles",
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
