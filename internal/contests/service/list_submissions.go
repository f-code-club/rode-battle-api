package service

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/f-code-club/rode-battle-api/internal/contests/repository"
	"github.com/f-code-club/rode-battle-api/internal/shared/errors"
	"github.com/google/uuid"
)

type Language = repository.Language

type Verdict = repository.Verdict

type SubmissionListItem struct {
	ID          uuid.UUID `json:"id"`
	ProblemName string    `json:"problem_name"`
	OwnerName   string    `json:"owner_name"`
	Language    Language  `json:"language"`
	Verdict     *Verdict  `json:"verdict"`
	Score       *float32  `json:"score"`
	CreatedAt   time.Time `json:"created_at"`
}

type ListSubmissionsFilter struct {
	ProblemID uuid.UUID
	AccountID uuid.UUID
	Language  *Language
	Verdict   *Verdict
	Search    *string
	Page      int
	PageSize  int
}

type ListSubmissionsResult struct {
	Items      []SubmissionListItem `json:"items"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	TotalItems int64                `json:"total_items"`
	TotalPages int                  `json:"total_pages"`
}

func (s *Service) ListSubmissionsByContest(
	ctx context.Context,
	contestID uuid.UUID,
	filter ListSubmissionsFilter,
) (*ListSubmissionsResult, error) {
	var search *string
	if filter.Search != nil {
		if trimmed := strings.TrimSpace(*filter.Search); trimmed != "" {
			search = &trimmed
		}
	}

	queries := repository.New(s.pool)

	rows, err := queries.GetSubmissionsByContest(ctx, repository.GetSubmissionsByContestParams{
		ContestID:  contestID,
		ProblemID:  filter.ProblemID,
		AccountID:  filter.AccountID,
		Language:   filter.Language,
		Verdict:    filter.Verdict,
		Search:     search,
		PageLimit:  int32(filter.PageSize),
		PageOffset: int32((filter.Page - 1) * filter.PageSize),
	})
	if err != nil {
		return nil, errors.Wrap(
			http.StatusInternalServerError,
			"failed to list submissions",
			err,
		)
	}

	var total int64
	if len(rows) > 0 {
		total = rows[0].TotalCount
	}

	items := make([]SubmissionListItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, SubmissionListItem{
			ID:          r.ID,
			ProblemName: r.ProblemName,
			OwnerName:   r.OwnerName,
			Language:    r.Language,
			Verdict:     r.Verdict,
			Score:       r.Score,
			CreatedAt:   r.CreatedAt,
		})
	}

	totalPages := int((total + int64(filter.PageSize) - 1) / int64(filter.PageSize))

	return &ListSubmissionsResult{
		Items:      items,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}
