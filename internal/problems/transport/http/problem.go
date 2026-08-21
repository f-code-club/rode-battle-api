package http

import (
	"github.com/go-fuego/fuego"
	"github.com/google/uuid"

	"github.com/f-code-club/rode-battle-api/internal/problems/service"
	"github.com/f-code-club/rode-battle-api/internal/shared/middleware"
)

func (s *Server) GetProblem(c fuego.ContextNoBody) (*service.Problem, error) {
	id := c.PathParam("id")

	problemId, err := uuid.Parse(id)
	if err != nil {
		return nil, fuego.BadRequestError{Title: "Invalid uuid", Err: err}
	}

	return s.service.GetProblem(c.Context(), problemId)
}

func (s *Server) GetSubmitHistory(c fuego.ContextNoBody) ([]service.ProblemHistory, error) {
	id := c.PathParam("id")
	accountId := c.Context().Value(middleware.AccountIDKey).(uuid.UUID)

	problemId, err := uuid.Parse(id)
	if err != nil {
		return nil, fuego.BadRequestError{Title: "Invalid uuid", Err: err}
	}

	return s.service.GetSubmitHistory(c.Context(), problemId, accountId)
}
