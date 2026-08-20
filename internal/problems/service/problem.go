package service

import (
	"context"
	"net/http"

	"github.com/f-code-club/rode-battle-api/internal/problems/repository"
	"github.com/f-code-club/rode-battle-api/internal/shared/errors"
	"github.com/google/uuid"
)

type Language = repository.Language

type Problem struct {
	Position    int32      `json:"position"`
	Name        string     `json:"name"`
	ContentPath string     `json:"content_path"`
	TimeLimit   int32      `json:"time_limit"`
	MemoryLimit int32      `json:"memory_limit"`
	Languages   []Language `json:"languages"`
}

func (s *Service) GetProblem(ctx context.Context, id uuid.UUID) (*Problem, error) {
	queries := repository.New(s.pool)

	problem, err := queries.GetProblem(ctx, id)
	if err != nil {
		return nil, errors.Wrap(http.StatusBadRequest, "Failed to get problem", err)
	}

	languages, err := queries.GetProblemLanguages(ctx, id)
	if err != nil {
		return nil, errors.Wrap(http.StatusBadRequest, "Failed to get language", err)
	}

	return &Problem{
		Position:    *problem.Position,
		Name:        problem.Name,
		ContentPath: problem.ContentPath,
		TimeLimit:   *problem.TimeLimit,
		MemoryLimit: *problem.MemoryLimit,
		Languages:   languages,
	}, nil
}
