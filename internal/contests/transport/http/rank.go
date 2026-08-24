package http

import (
	"net/http"

	"github.com/f-code-club/rode-battle-api/internal/contests/service"
	"github.com/f-code-club/rode-battle-api/internal/shared/errors"
	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
)

func (s *Server) GetRank(c fuego.ContextNoBody) ([]service.Ranking, error) {
	id := c.PathParam("id")

	contestID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.Wrap(
			http.StatusBadRequest,
			"invalid uuid",
			err,
		)
	}

	return s.service.GetRank(c.Context(), contestID)
}
