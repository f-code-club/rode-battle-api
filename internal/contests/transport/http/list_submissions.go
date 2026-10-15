package http

import (
	"context"

	"github.com/f-code-club/rode-battle-api/internal/contests/service"
	"github.com/f-code-club/rode-battle-api/internal/shared/pagination"
	"github.com/google/uuid"
)

type ListSubmissionsByContestInput struct {
	pagination.Params
	ContestID uuid.UUID `path:"id"`
	ProblemID uuid.UUID `query:"problem_id"`
	AccountID uuid.UUID `query:"account_id"`
	Language  string    `query:"language" enum:"rust,cpp,python,java,html"`
	Verdict   string    `query:"verdict" enum:"accepted,wrong_answer,time_limit_exceeded,compilation_error,memory_limit_exceeded,runtime_error,idle_time_limit_exceeded"`
	Search    string    `query:"search" maxLength:"100"`
}

type ListSubmissionsByContestOutput struct {
	Body *pagination.Response[service.SubmissionListItem]
}

func (s *Server) ListSubmissionsByContest(
	ctx context.Context,
	input *ListSubmissionsByContestInput,
) (*ListSubmissionsByContestOutput, error) {
	filter := service.ListSubmissionsFilter{
		Params:    input.Params,
		ProblemID: input.ProblemID,
		AccountID: input.AccountID,
	}

	if input.Language != "" {
		l := service.Language(input.Language)
		filter.Language = &l
	}
	if input.Verdict != "" {
		v := service.Verdict(input.Verdict)
		filter.Verdict = &v
	}
	if input.Search != "" {
		filter.Search = &input.Search
	}

	result, err := s.service.ListSubmissionsByContest(ctx, input.ContestID, filter)
	if err != nil {
		return nil, err
	}
	return &ListSubmissionsByContestOutput{Body: result}, nil
}
