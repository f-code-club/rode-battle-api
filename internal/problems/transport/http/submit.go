package http

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/f-code-club/rode-battle-api/internal/problems/repository"
	"github.com/f-code-club/rode-battle-api/internal/shared/middleware"
)

type Language = repository.Language

type SubmissionRequest struct {
	Language Language `json:"language" required:"true"`
	Code     string   `json:"code" required:"true"`
}

type CreateSubmissionInput struct {
	ID   uuid.UUID `path:"id"`
	Body SubmissionRequest
}

type ResentSubmissionInput struct {
	ID uuid.UUID `path:"id"`
}

type CreateSubmissionOutput struct {
	Body uuid.UUID
}

func (s *Server) CreateSubmission(ctx context.Context, input *CreateSubmissionInput) (*CreateSubmissionOutput, error) {
	accountID, ok := ctx.Value(middleware.AccountIDKey).(uuid.UUID)
	if !ok {
		return nil, huma.Error401Unauthorized("unauthorized")
	}

	id, err := s.service.CreateSubmission(ctx, input.ID, accountID, input.Body.Language, input.Body.Code)
	if err != nil {
		return nil, err
	}

	return &CreateSubmissionOutput{Body: id}, nil
}
