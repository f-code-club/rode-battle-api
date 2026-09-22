package http

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/f-code-club/rode-battle-api/internal/auth/service"
	"github.com/f-code-club/rode-battle-api/internal/shared/middleware"
)

type MeInput struct{}

type MeOutput struct {
	Body *service.Account
}

func (s *Server) Me(ctx context.Context, input *MeInput) (*MeOutput, error) {
	id, ok := ctx.Value(middleware.AccountIDKey).(uuid.UUID)
	if !ok {
		return nil, huma.Error401Unauthorized("unauthorized")
	}

	acc, err := s.service.Me(ctx, id)
	if err != nil {
		return nil, err
	}

	return &MeOutput{Body: acc}, nil
}
