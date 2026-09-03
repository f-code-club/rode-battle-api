package http

import (
	"github.com/f-code-club/rode-battle-api/internal/contests/service"
	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
)

func (s *Server) CreateContest(c fuego.ContextWithBody[service.CreateContestRequest]) (uuid.UUID, error) {
	req, err := c.Body()
	if err != nil {
		return uuid.Nil, fuego.BadRequestError{
			Title: "Invalid request body",
			Err:   err,
		}
	}

	return s.service.CreateContest(c.Context(), req)
}
