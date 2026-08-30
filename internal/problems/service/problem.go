package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/f-code-club/rode-battle-api/internal/problems/repository"
	apperr "github.com/f-code-club/rode-battle-api/internal/shared/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type Language = repository.Language

type Verdict = repository.Verdict

type GetSubmitHistory = repository.GetSubmitHistoryParams

type CreateProblem = repository.CreateProblemParams

type CreateProblemLanguage = repository.CreateProblemLanguageParams

var algorithmLanguages = map[string]struct{}{
	"rust":   {},
	"cpp":    {},
	"python": {},
	"java":   {},
}

type Problem struct {
	Position    *int32     `json:"position"`
	Name        string     `json:"name"`
	Content     string     `json:"content"`
	TimeLimit   *int32     `json:"time_limit"`
	MemoryLimit *int32     `json:"memory_limit"`
	Languages   []Language `json:"languages"`
}

type ProblemHistory struct {
	ID        uuid.UUID `json:"id"`
	Language  Language  `json:"language"`
	Code      string    `json:"code"`
	Verdict   *Verdict  `json:"verdict"`
	Score     *float32  `json:"score"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateProblemInput struct {
	Name            string
	Content         string
	CheckerLanguage *Language
	CheckerPath     *string
	TimeLimit       *int32
	MemoryLimit     *int32
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
		Position:    problem.Position,
		Name:        problem.Name,
		Content:     problem.Content,
		TimeLimit:   problem.TimeLimit,
		MemoryLimit: problem.MemoryLimit,
		Languages:   languages,
	}, nil
}

func (s *Service) GetSubmitHistory(ctx context.Context, problemID uuid.UUID, accountID uuid.UUID) ([]ProblemHistory, error) {
	queries := repository.New(s.pool)

	rows, err := queries.GetSubmitHistory(ctx, GetSubmitHistory{
		ProblemID: problemID,
		AccountID: accountID,
	})
	if err != nil {
		return nil, apperr.Wrap(http.StatusBadRequest, "Problem not found", err)
	}

	history := make([]ProblemHistory, 0, len(rows))
	for _, row := range rows {
		history = append(history, ProblemHistory{
			ID:        row.ID,
			Language:  row.Language,
			Code:      row.Code,
			Verdict:   row.Verdict,
			Score:     row.Score,
			CreatedAt: row.CreatedAt,
		})
	}

	return history, nil
}

func (s *Service) CreateProblem(ctx context.Context, input CreateProblemInput, language []string) (uuid.UUID, error) {
	var pgErr *pgconn.PgError
	requiredAlgoInput := false
	for _, lang := range language {
		if _, ok := algorithmLanguages[lang]; ok {
			requiredAlgoInput = true
		}

		if len(language) > 1 && lang == "html" {
			return uuid.Nil, apperr.Wrap(http.StatusBadRequest, "Language mismatch", nil)
		}
	}

	if requiredAlgoInput && (input.CheckerPath == nil || input.CheckerLanguage == nil || input.MemoryLimit == nil || input.TimeLimit == nil) {
		return uuid.Nil, apperr.Wrap(http.StatusBadRequest, "Cannot leave checker_code, checker_language, time_limit, memory_limit empty", nil)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, apperr.Wrap(http.StatusInternalServerError, "Failed to create problem", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	qtx := repository.New(s.pool).WithTx(tx)

	rows, err := qtx.CreateProblem(ctx, CreateProblem{
		Name:            input.Name,
		Content:         input.Content,
		CheckerLanguage: input.CheckerLanguage,
		CheckerPath:     input.CheckerPath,
		TimeLimit:       input.TimeLimit,
		MemoryLimit:     input.MemoryLimit,
	})
	if err != nil {
		return uuid.Nil, apperr.Wrap(http.StatusBadRequest, "Failed to create problem", err)
	}

	err = qtx.CreateProblemLanguage(ctx, CreateProblemLanguage{
		ProblemID: rows,
		Language:  language,
	})
	if ok := errors.As(err, &pgErr); ok {
		switch pgErr.Code {
		case "23505":
			return uuid.Nil, apperr.Wrap(http.StatusBadRequest, "Duplicate language", err)
		case "22P02":
			return uuid.Nil, apperr.Wrap(http.StatusBadRequest, "Invalid language", err)
		}
	}
	if err != nil {
		return uuid.Nil, apperr.Wrap(http.StatusBadRequest, "Failed to create problem language", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, apperr.Wrap(http.StatusInternalServerError, "Failed to create problem", err)
	}

	return rows, nil

}
