package http

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

type RefreshInput struct {
	RefreshToken string `cookie:"refresh_token"`
}

type RefreshOutput struct {
	Body string
}

func (s *Server) Refresh(ctx context.Context, input *RefreshInput) (*RefreshOutput, error) {
	if input.RefreshToken == "" {
		return nil, huma.Error401Unauthorized("refresh token not found")
	}

	accessToken, err := s.service.Refresh(ctx, input.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &RefreshOutput{Body: accessToken}, nil
}
