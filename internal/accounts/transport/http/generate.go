package http

import (
	"context"

	"github.com/google/uuid"

	"github.com/f-code-club/rode-battle-api/internal/accounts/service"
)

type Role string

const (
	Participant Role = "participant"
	Jury        Role = "jury"
	Admin       Role = "admin"
)

type GenerateRequest struct {
	Email string `json:"email" format:"email" required:"true"`
	Name  string `json:"name" required:"true"`
	Role  Role   `json:"role" enum:"participant,jury,admin" required:"true"`
}

type GenerateInput struct {
	Body GenerateRequest
}

type GenerateOutput struct {
	Body uuid.UUID
}

func (s *Server) Generate(ctx context.Context, input *GenerateInput) (*GenerateOutput, error) {
	id, err := s.service.Generate(ctx, input.Body.Email, input.Body.Name, service.Role(input.Body.Role))
	if err != nil {
		return nil, err
	}

	return &GenerateOutput{Body: id}, nil
}
