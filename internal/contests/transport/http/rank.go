package http

import (
	"context"

	"github.com/google/uuid"

	"github.com/f-code-club/rode-battle-api/internal/contests/service"
)

type GetRankInput struct {
	ID uuid.UUID `path:"id"`
}

type GetRankOutput struct {
	Body []service.Ranking
}

func (s *Server) GetRank(ctx context.Context, input *GetRankInput) (*GetRankOutput, error) {
	rankings, err := s.service.GetRank(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	return &GetRankOutput{Body: rankings}, nil
}
