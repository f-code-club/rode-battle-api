package service

import (
	"context"
	"net/http"

	"github.com/f-code-club/rode-battle-api/internal/problems/repository"
	apperr "github.com/f-code-club/rode-battle-api/internal/shared/errors"
	"github.com/google/uuid"
)

type CreateTestCaseParams = repository.CreateTestCaseParams

func (s *Service) CreateTestCase(ctx context.Context, problemID uuid.UUID, input string) (uuid.UUID, error) {
	queries := repository.New(s.pool)
	if input == "" {
		return uuid.Nil, apperr.Wrap(http.StatusBadRequest, "Empty input_path", nil)
	}

	testCaseID, err := queries.CreateTestCase(ctx, CreateTestCaseParams{
		ProblemID: problemID,
		Input:     input,
	})
	if err != nil {
		return uuid.Nil, apperr.Wrap(http.StatusInternalServerError, "Failed to create test case", err)
	}

	return testCaseID, nil
}
