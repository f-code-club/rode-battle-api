package service

import (
	"context"
	"net/http"
	"time"

	"github.com/f-code-club/rode-battle-api/internal/contest/repository"
	"github.com/f-code-club/rode-battle-api/internal/shared/errors"
	"github.com/google/uuid"
)

type Contest struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func (s *Service) ListContests(
	ctx context.Context,
) ([]Contest, error) {
	queries := repository.New(s.pool)

	contests, err := queries.GetContests(ctx)
	if err != nil {
		return nil, errors.Wrap(
			http.StatusInternalServerError,
			"failed to list contests",
			err,
		)
	}

	result := make([]Contest, 0, len(contests))
	for _, c := range contests {
		result = append(result, Contest{
			ID:    c.ID,
			Name:  c.Name,
			Start: c.Start.Time,
			End:   c.End.Time,
		})
	}

	return result, nil
}
