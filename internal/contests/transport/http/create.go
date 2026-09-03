package http

import (
	"github.com/f-code-club/rode-battle-api/internal/contests/service"
	"github.com/go-fuego/fuego"
)

func (s *Server) CreateContest(c fuego.ContextWithBody[service.CreateContestRequest]) (string, error) {
	req, err := c.Body()
	if err != nil {
		return "", fuego.BadRequestError{
			Title: "Invalid request body",
			Err:   err,
		}
	}

	contestID, err := s.service.CreateContest(c.Context(), req)
	if err != nil {
		return "", err
	}

	return contestID.String(), nil
}
