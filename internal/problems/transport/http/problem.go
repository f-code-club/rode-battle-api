package http

import (
	"github.com/go-fuego/fuego"
	"github.com/google/uuid"

	"github.com/f-code-club/rode-battle-api/internal/problems/service"
	"github.com/f-code-club/rode-battle-api/internal/shared/middleware"
)

type CreateProblemRequest struct {
	Name            string    `json:"name" validate:"required"`
	Content         string    `json:"content" validate:"required"`
	CheckerLanguage *Language `json:"checker_language"`
	CheckerCode     *string   `json:"checker_code"`
	TimeLimit       *int32    `json:"time_limit"`
	MemoryLimit     *int32    `json:"memory_limit"`
	Languages       []string  `json:"languages" validate:"required,unique,min=1,dive,oneof=rust cpp python java html"`
}

func (s *Server) GetProblem(c fuego.ContextNoBody) (*service.Problem, error) {
	id := c.PathParam("id")

	problemID, err := uuid.Parse(id)
	if err != nil {
		return nil, fuego.BadRequestError{Title: "Invalid uuid", Err: err}
	}

	return s.service.GetProblem(c.Context(), problemID)
}

func (s *Server) GetSubmitHistory(c fuego.ContextNoBody) ([]service.ProblemHistory, error) {
	id := c.PathParam("id")
	accountID := c.Context().Value(middleware.AccountIDKey).(uuid.UUID)

	problemID, err := uuid.Parse(id)
	if err != nil {
		return nil, fuego.BadRequestError{Title: "Invalid uuid", Err: err}
	}

	return s.service.GetSubmitHistory(c.Context(), problemID, accountID)
}

func (s *Server) CreateProblem(c fuego.ContextWithBody[CreateProblemRequest]) (uuid.UUID, error) {
	body, err := c.Body()
	if err != nil {
		return uuid.Nil, err
	}

	return s.service.CreateProblem(c, service.CreateProblemInput{
		Name:            body.Name,
		Content:         body.Content,
		CheckerLanguage: body.CheckerLanguage,
		CheckerPath:     body.CheckerCode,
		TimeLimit:       body.TimeLimit,
		MemoryLimit:     body.MemoryLimit,
	}, body.Languages)
}
