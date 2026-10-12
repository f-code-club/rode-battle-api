package http

import (
	"context"

	"github.com/f-code-club/rode-battle-api/internal/contests/service"
	"github.com/google/uuid"
)

type ListSubmissionsByContestInput struct {
	ContestID uuid.UUID `path:"id"`

	ProblemID uuid.UUID `query:"problem_id"`
	AccountID uuid.UUID `query:"account_id"`
	Language  string    `query:"language" enum:"rust,cpp,python,java,html"`
	Verdict   string    `query:"verdict" enum:"accepted,wrong_answer,time_limit_exceeded,compilation_error,memory_limit_exceeded,runtime_error,idle_time_limit_exceeded"`
	Search    string    `query:"search" maxLength:"100"`

	Page     int `query:"page" default:"1" minimum:"1" maximum:"100000"`
	PageSize int `query:"page_size" default:"20" minimum:"1" maximum:"100"`
}

type ListSubmissionsByContestOutput struct {
	Body *service.ListSubmissionsResult
}

func (s *Server) ListSubmissionsByContest(
	ctx context.Context,
	input *ListSubmissionsByContestInput,
) (*ListSubmissionsByContestOutput, error) {
	filter := service.ListSubmissionsFilter{
		ProblemID: input.ProblemID,
		AccountID: input.AccountID,
		Page:      input.Page,
		PageSize:  input.PageSize,
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
