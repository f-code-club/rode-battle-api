package http

import (
	"github.com/go-fuego/fuego"
	"github.com/google/uuid"

	"github.com/f-code-club/rode-battle-api/internal/auth/service"
	"github.com/f-code-club/rode-battle-api/internal/shared/middleware"
)

func (s *Server) Me(c fuego.ContextNoBody) (*service.Account, error) {
	id := c.Context().Value(middleware.AccountIDKey).(uuid.UUID)

	return s.service.Me(c.Context(), id)
}
