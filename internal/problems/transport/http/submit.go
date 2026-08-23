package http

import (
	"github.com/f-code-club/rode-battle-api/internal/problems/repository"
	"github.com/f-code-club/rode-battle-api/internal/shared/middleware"
	"github.com/go-fuego/fuego"
	"github.com/google/uuid"
)

type Language = repository.Language

type SubmissionRequest struct {
	Language Language `json:"language" validate:"required"`
	Code     string   `json:"code" validate:"required"`
}

func (s *Server) CreateSubmission(c fuego.ContextWithBody[SubmissionRequest]) (uuid.UUID, error) {
	body, err := c.Body()
	if err != nil {
		return uuid.Nil, err
	}
	id := c.PathParam("id")
	accountID := c.Context().Value(middleware.AccountIDKey).(uuid.UUID)

	problemId, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, err
	}

	return s.service.CreateSubmission(c.Context(), problemId, accountID, body.Language, body.Code)
}
