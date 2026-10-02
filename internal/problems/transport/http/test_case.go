package http

import (
	"context"

	"github.com/google/uuid"
)

type CreateTestCaseRequest struct {
	InputPath string `json:"input_path" required:"true"`
}

type CreateTestCaseInput struct {
	ID   uuid.UUID `path:"id"`
	Body CreateTestCaseRequest
}

type CreateTestCaseOutput struct {
	Body uuid.UUID
}

func (s *Server) CreateTestCase(ctx context.Context, input *CreateTestCaseInput) (*CreateTestCaseOutput, error) {
	id, err := s.service.CreateTestCase(ctx, input.ID, input.Body.InputPath)
	if err != nil {
		return nil, err
	}

	return &CreateTestCaseOutput{Body: id}, nil
}
