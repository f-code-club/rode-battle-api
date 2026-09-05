package service

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/f-code-club/rode-battle-api/internal/contests/repository"
	"github.com/f-code-club/rode-battle-api/internal/shared/errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var validate = validator.New()

type CreateContestRequest struct {
	Name     string      `json:"name" validate:"required"`
	Start    time.Time   `json:"start" validate:"required"`
	End      time.Time   `json:"end" validate:"required,gtfield=Start"`
	Problems []uuid.UUID `json:"problems"`
}

func (s *Service) CreateContest(
	ctx context.Context,
	req CreateContestRequest,
) (uuid.UUID, error) {
	req.Name = strings.TrimSpace(req.Name)

	if err := validate.Struct(req); err != nil {
		return uuid.Nil, errors.Wrap(
			http.StatusBadRequest,
			"invalid contest request",
			err,
		)
	}

	if err := validateContestProblems(req.Problems); err != nil {
		return uuid.Nil, err
	}

	queries := repository.New(s.pool)

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, errors.Wrap(
			http.StatusInternalServerError,
			"failed to begin transaction",
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
			"failed to validate contest problems",
			err,
		)
	}

	if err := validateProblemAssignment(problems, req.Problems); err != nil {
		return uuid.Nil, err
	}

	contestID, err := txQueries.CreateContest(ctx, repository.CreateContestParams{
		Name: req.Name,
		StartTime: pgtype.Timestamptz{
			Time:  req.Start,
			Valid: true,
		},
		EndTime: pgtype.Timestamptz{
			Time:  req.End,
			Valid: true,
		},
	})
	if err != nil {
		return uuid.Nil, errors.Wrap(
			http.StatusInternalServerError,
			"failed to create contest",
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
				"failed to assign problems to contest",
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
			"failed to commit contest creation",
			err,
		)
	}

	return contestID, nil
}

func validateContestProblems(problemIDs []uuid.UUID) error {
	seen := make(map[uuid.UUID]struct{}, len(problemIDs))

	for _, problemID := range problemIDs {
		if problemID == uuid.Nil {
			return errors.Wrap(
				http.StatusBadRequest,
				"problem id must not be empty",
				nil,
			)
		}

		if _, exists := seen[problemID]; exists {
			return errors.Wrap(
				http.StatusBadRequest,
				"duplicate problem id",
				nil,
			)
		}

		seen[problemID] = struct{}{}
	}

	return nil
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
