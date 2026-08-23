package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/f-code-club/rode-battle-api/internal/problems/repository"
	apperr "github.com/f-code-club/rode-battle-api/internal/shared/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type CreateSubmission = repository.CreateSubmissionParams

func (s *Service) CreateSubmission(ctx context.Context, problemID uuid.UUID, accountID uuid.UUID, language Language, code string) (uuid.UUID, error) {
	var pgErr *pgconn.PgError
	queries := repository.New(s.pool)

	submissionID, err := queries.CreateSubmission(ctx, CreateSubmission{
		ProblemID: problemID,
		AccountID: accountID,
		Language:  language,
		Code:      code,
	})
	if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
		return uuid.Nil, apperr.Wrap(http.StatusBadRequest, "Invalid language", err)
	}
	if err != nil {
		return uuid.Nil, apperr.Wrap(http.StatusBadRequest, "Failed to create submission", err)
	}
	// Add queue later

	return submissionID, nil
}
