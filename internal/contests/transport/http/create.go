package http

import (
	"time"

	"github.com/f-code-club/rode-battle-api/internal/contests/service"
	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
)

type CreateContestRequest struct {
	Name     string      `json:"name" validate:"required"`
	Start    time.Time   `json:"start" validate:"required"`
	End      time.Time   `json:"end" validate:"required,gtfield=Start"`
	Problems []uuid.UUID `json:"problems" validate:"unique,dive,required"`
}

func (s *Server) CreateContest(c fuego.ContextWithBody[CreateContestRequest]) (uuid.UUID, error) {
	body, err := c.Body()
	if err != nil {
		return uuid.Nil, fuego.BadRequestError{
			Title: "Invalid request body",
			Err:   err,
		}
	}

	return s.service.CreateContest(c.Context(), service.CreateContestInput{
		Name:     body.Name,
		Start:    body.Start,
		End:      body.End,
		Problems: body.Problems,
	})
}
