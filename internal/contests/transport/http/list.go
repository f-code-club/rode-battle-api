package http

import (
	"context"

	"github.com/f-code-club/rode-battle-api/internal/contests/service"
)

type ListContestsInput struct{}

type ListContestsOutput struct {
	Body []service.Contest
}

func (s *Server) ListContests(ctx context.Context, input *ListContestsInput) (*ListContestsOutput, error) {
	contests, err := s.service.ListContests(ctx)
	if err != nil {
		return nil, err
	}
	return &ListContestsOutput{Body: contests}, nil
}
