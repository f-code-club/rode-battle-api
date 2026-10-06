package service

import (
	"context"
	"net/http"
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

func (s *Service) ListSubmissionsByContest(
	ctx context.Context,
	contestID uuid.UUID,
) ([]SubmissionListItem, error) {
	queries := repository.New(s.pool)

	rows, err := queries.GetSubmissionsByContest(ctx, contestID)
	if err != nil {
		return nil, errors.Wrap(
			http.StatusInternalServerError,
			"failed to list submissions",
			err,
		)
	}

	result := make([]SubmissionListItem, 0, len(rows))
	for _, r := range rows {
		result = append(result, SubmissionListItem{
			ID:          r.ID,
			ProblemName: r.ProblemName,
			OwnerName:   r.OwnerName,
			Language:    r.Language,
			Verdict:     r.Verdict,
			Score:       r.Score,
			CreatedAt:   r.CreatedAt,
		})
	}

	return result, nil
}
