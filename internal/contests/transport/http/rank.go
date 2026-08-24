package http

import (
	"github.com/f-code-club/rode-battle-api/internal/contests/service"
	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
)

func (s *Server) GetRank(c fuego.ContextNoBody) ([]service.Ranking, error) {
	id := c.PathParam("id")

	contestID, err := uuid.Parse(id)
	if err != nil {
		return nil, fuego.BadRequestError{Title: "Invalid uuid", Err: err}
	}

	return s.service.GetRank(c.Context(), contestID)
}
