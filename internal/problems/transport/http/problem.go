package http

import (
	"github.com/go-fuego/fuego"
	"github.com/google/uuid"

	"github.com/f-code-club/rode-battle-api/internal/problems/service"
)

func (s *Server) GetProblem(c fuego.ContextNoBody) (*service.Problem, error) {
	id := c.PathParam("id")

	problemId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return s.service.GetProblem(c.Context(), problemId)
}
