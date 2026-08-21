package service

import (
	"context"
	"net/http"
	"time"

	"github.com/f-code-club/rode-battle-api/internal/problems/repository"
	"github.com/f-code-club/rode-battle-api/internal/shared"
	apperr "github.com/f-code-club/rode-battle-api/internal/shared/errors"
	"github.com/google/uuid"
)

type Language = repository.Language

type Verdict = repository.Verdict

type Problem struct {
	Position    int32      `json:"position"`
	Name        string     `json:"name"`
	ContentPath string     `json:"content_path"`
	TimeLimit   int32      `json:"time_limit"`
	MemoryLimit int32      `json:"memory_limit"`
	Languages   []Language `json:"languages"`
}

type ProblemHistory struct {
	ID        uuid.UUID `json:"id"`
	Language  Language  `json:"language"`
	Code      string    `json:"code"`
	Verdict   Verdict   `json:"verdict"`
	Score     float32   `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Service) GetProblem(ctx context.Context, id uuid.UUID) (*Problem, error) {
	queries := repository.New(s.pool)

	problem, err := queries.GetProblem(ctx, id)
	if err != nil {
		return nil, apperr.Wrap(http.StatusBadRequest, "Failed to get problem", err)
	}

	languages, err := queries.GetProblemLanguages(ctx, id)
	if err != nil {
		return nil, apperr.Wrap(http.StatusBadRequest, "Failed to get language", err)
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

func (s *Service) GetSubmitHistory(ctx context.Context, problemID uuid.UUID) ([]ProblemHistory, error) {
	queries := repository.New(s.pool)

	rows, err := queries.GetSubmitHistory(ctx, problemID)
	if err != nil {
		return nil, apperr.Wrap(http.StatusBadRequest, "Problem not found", err)
	}

	history := make([]ProblemHistory, 0, len(rows))
	println(len(history))
	for _, row := range rows {
		history = append(history, ProblemHistory{
			ID:        row.ID,
			Language:  row.Language,
			Code:      row.Code,
			Verdict:   shared.Deref(row.Verdict),
			Score:     shared.Deref(row.Score),
			CreatedAt: row.CreatedAt,
		})
	}

	return history, nil
}
