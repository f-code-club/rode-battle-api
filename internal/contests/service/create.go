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

const internalErrorMessage = "something went wrong"

func (s *Service) CreateContest(
	ctx context.Context,
	name string,
	start, end time.Time,
	problems []uuid.UUID,
) (uuid.UUID, error) {
	if len(problems) == 0 {
		return uuid.Nil, errors.New(http.StatusBadRequest, "no problem assigned to contest")
	}

	name = strings.TrimSpace(name)

	queries := repository.New(s.pool)

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, errors.Wrap(
			http.StatusInternalServerError,
			internalErrorMessage,
			err,
		)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	txQueries := queries.WithTx(tx)

	contestID, err := txQueries.CreateContest(ctx, repository.CreateContestParams{
		Name:      name,
		StartTime: start,
		EndTime:   end,
	})
	if err != nil {
		return uuid.Nil, errors.Wrap(
			http.StatusInternalServerError,
			internalErrorMessage,
			err,
		)
	}
	for i, id := range problems {
		position := int32(i)

		err := txQueries.AssignProblemToContest(
			ctx,
			repository.AssignProblemToContestParams{
				ContestID: contestID,
				ID:        id,
				Position:  &position,
			},
		)
		if err != nil {
			return uuid.Nil, errors.Wrap(
				http.StatusInternalServerError,
				internalErrorMessage,
				err,
			)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, errors.Wrap(
			http.StatusInternalServerError,
			internalErrorMessage,
			err,
		)
	}

	return contestID, nil
}
