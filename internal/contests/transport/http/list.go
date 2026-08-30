package http

import (
	"github.com/f-code-club/rode-battle-api/internal/contests/service"
	"github.com/go-fuego/fuego"
)

func (s *Server) ListContests(c fuego.ContextNoBody) ([]service.Contest, error) {
	return s.service.ListContests(c.Context())
}
