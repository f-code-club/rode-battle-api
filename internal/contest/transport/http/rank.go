package http

import (
	"fmt"

	"github.com/f-code-club/rode-battle-api/internal/contest/service"
	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
)

func (s *Server) GetRank(c fuego.ContextNoBody) ([]service.Ranking, error) {
	id := c.PathParam("id")

	fmt.Println(id)

	contestID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return s.service.GetRank(c.Context(), contestID)
}
