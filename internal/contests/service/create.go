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
	problemIds []uuid.UUID,
) (uuid.UUID, error) {
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

	problems, err := txQueries.GetProblemsForContestAssignment(ctx, problemIds)
	if err != nil {
		return uuid.Nil, errors.Wrap(
			http.StatusInternalServerError,
			internalErrorMessage,
			err,
		)
	}

	if err := validateProblemAssignment(problems, problemIds); err != nil {
		return uuid.Nil, err
	}

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

	if len(problemIds) > 0 {
		assigned, err := txQueries.AssignProblemsToContest(
			ctx,
			repository.AssignProblemsToContestParams{
				ContestID:  contestID,
				ProblemIds: problemIds,
			},
		)
		if err != nil {
			return uuid.Nil, errors.Wrap(
				http.StatusInternalServerError,
				internalErrorMessage,
				err,
			)
		}

		if assigned != int64(len(problemIds)) {
			return uuid.Nil, errors.Wrap(
				http.StatusConflict,
				"one or more problems are already assigned to another contest",
				nil,
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

func validateProblemAssignment(
	problems []repository.GetProblemsForContestAssignmentRow,
	requestedIDs []uuid.UUID,
) error {
	if len(problems) != len(requestedIDs) {
		return errors.Wrap(
			http.StatusBadRequest,
			"one or more problems do not exist",
			nil,
		)
	}

	for _, problem := range problems {
		if problem.ContestID != uuid.Nil {
			return errors.Wrap(
				http.StatusConflict,
				"one or more problems are already assigned to another contest",
				nil,
			)
		}
	}

	return nil
}
