package http

import (
	"context"

	"github.com/f-code-club/rode-battle-api/internal/contests/service"
	"github.com/google/uuid"
)

type ListSubmissionsByContestInput struct {
	ContestID uuid.UUID `path:"id"`
}

type ListSubmissionsByContestOutput struct {
	Body []service.SubmissionListItem
}

func (s *Server) ListSubmissionsByContest(
	ctx context.Context,
	input *ListSubmissionsByContestInput,
) (*ListSubmissionsByContestOutput, error) {
	submissions, err := s.service.ListSubmissionsByContest(ctx, input.ContestID)
	if err != nil {
		return nil, err
	}
	return &ListSubmissionsByContestOutput{Body: submissions}, nil
}
