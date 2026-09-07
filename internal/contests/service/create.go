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

type CreateContestInput struct {
	Name     string
	Start    time.Time
	End      time.Time
	Problems []uuid.UUID
}

func (s *Service) CreateContest(
	ctx context.Context,
	req CreateContestInput,
) (uuid.UUID, error) {
	req.Name = strings.TrimSpace(req.Name)

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

	problems, err := txQueries.GetProblemsForContestAssignment(ctx, req.Problems)
	if err != nil {
		return uuid.Nil, errors.Wrap(
			http.StatusInternalServerError,
			internalErrorMessage,
			err,
		)
	}

	if err := validateProblemAssignment(problems, req.Problems); err != nil {
		return uuid.Nil, err
	}

	contestID, err := txQueries.CreateContest(ctx, repository.CreateContestParams{
		Name:      req.Name,
		StartTime: req.Start,
		EndTime:   req.End,
	})
	if err != nil {
		return uuid.Nil, errors.Wrap(
			http.StatusInternalServerError,
			internalErrorMessage,
			err,
		)
	}

	if len(req.Problems) > 0 {
		assigned, err := txQueries.AssignProblemsToContest(
			ctx,
			repository.AssignProblemsToContestParams{
				ContestID:  contestID,
				ProblemIds: req.Problems,
			},
		)
		if err != nil {
			return uuid.Nil, errors.Wrap(
				http.StatusInternalServerError,
				internalErrorMessage,
				err,
			)
		}

		if assigned != int64(len(req.Problems)) {
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
