package http

import (
	"context"

	"github.com/google/uuid"

	"github.com/f-code-club/rode-battle-api/internal/contests/service"
)

type GetContestDetailInput struct {
	ID uuid.UUID `path:"id"`
}

type GetContestDetailOutput struct {
	Body service.ContestDetail
}

func (s *Server) GetContestDetail(ctx context.Context, input *GetContestDetailInput) (*GetContestDetailOutput, error) {
	detail, err := s.service.GetContestDetail(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return &GetContestDetailOutput{Body: detail}, nil
}
