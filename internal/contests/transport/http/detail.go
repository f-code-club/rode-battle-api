package http

import (
	"github.com/f-code-club/rode-battle-api/internal/contests/service"
	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
)

func (s *Server) GetContestDetail(c fuego.ContextNoBody) (service.ContestDetail, error) {
	id := c.PathParam("id")

	contestID, err := uuid.Parse(id)
	if err != nil {
		return service.ContestDetail{}, fuego.BadRequestError{Title: "Invalid uuid", Err: err}
	}

	return s.service.GetContestDetail(c.Context(), contestID)
}
