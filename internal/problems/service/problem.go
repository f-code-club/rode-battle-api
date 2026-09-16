package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/f-code-club/rode-battle-api/internal/problems/repository"
	"github.com/f-code-club/rode-battle-api/internal/shared"
	apperr "github.com/f-code-club/rode-battle-api/internal/shared/errors"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
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
	ColorCode   *string    `json:"color_code"`
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
	CheckerCode     *string
	TimeLimit       *int32
	MemoryLimit     *int32
	ColorCode       *string
}

type CompileRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

func (s *Service) GetProblem(ctx context.Context, id uuid.UUID) (*Problem, error) {
	queries := repository.New(s.pool)

	problem, err := queries.GetProblem(ctx, id)
	if err != nil {
		return nil, apperr.Wrap(http.StatusBadRequest, "Failed to get problem", err)
	}
	content := problem.Content

	languages, err := queries.GetProblemLanguages(ctx, id)
	if err != nil {
		return nil, apperr.Wrap(http.StatusBadRequest, "Failed to get language", err)
	}

	if languages[0] == "html" {
		content, err = s.s3.GetPresignedURl(ctx, problem.Content, 15*time.Minute)
		if err != nil {
			return nil, apperr.Wrap(http.StatusInternalServerError, "Failed to get get URL", err)
		}
	}

	return &Problem{
		Position:    problem.Position,
		Name:        problem.Name,
		Content:     content,
		TimeLimit:   problem.TimeLimit,
		MemoryLimit: problem.MemoryLimit,
		ColorCode:   problem.ColorCode,
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
	key, err := shared.RandomKey("problems")
	if err != nil {
		return uuid.Nil, apperr.Wrap(http.StatusInternalServerError, "Failed to random key", err)
	}

	requiredAlgoInput := false
	for _, lang := range language {
		if _, ok := algorithmLanguages[lang]; ok {
			requiredAlgoInput = true
		}

		if len(language) > 1 && lang == "html" {
			return uuid.Nil, apperr.Wrap(http.StatusBadRequest, "Language mismatch", nil)
		}
	}

	if requiredAlgoInput && (input.CheckerCode == nil || input.CheckerLanguage == nil || input.MemoryLimit == nil || input.TimeLimit == nil) {
		return uuid.Nil, apperr.Wrap(http.StatusBadRequest, "Cannot leave checker_code, checker_language, time_limit, memory_limit empty", nil)
	}

	content := input.Content
	var checkerPath *string
	if !requiredAlgoInput {
		contentKey, err := shared.RandomKey("problems")
		if err != nil {
			return uuid.Nil, apperr.Wrap(http.StatusInternalServerError, "Failed to random key", err)
		}
		decoded, err := base64.StdEncoding.DecodeString(input.Content)
		if err != nil {
			return uuid.Nil, apperr.Wrap(http.StatusInternalServerError, "Faield to decode base64", err)
		}
		mime := mimetype.Detect(decoded)
		file := bytes.NewReader(decoded)

		err = s.s3.UploadFile(ctx, contentKey, file, mime.String())
		if err != nil {
			return uuid.Nil, apperr.Wrap(http.StatusInternalServerError, "Failed to upload problems to storage", err)
		}

		content = contentKey
	} else {
		checkerBody, err := json.Marshal(CompileRequest{
			Code:     *input.CheckerCode,
			Language: string(*input.CheckerLanguage),
		})
		if err != nil {
			return uuid.Nil, apperr.Wrap(http.StatusInternalServerError, "Failed to create problem", err)
		}
		res, err := http.Post(fmt.Sprintf("%s/compile", s.judgeURL), "application/json; charset=utf-8", bytes.NewBuffer(checkerBody))
		if err != nil {
			return uuid.Nil, apperr.Wrap(http.StatusInternalServerError, "Failed to create problem", err)
		}
		defer func() {
			_ = res.Body.Close()
		}()
		if res.StatusCode != http.StatusOK {
			return uuid.Nil, apperr.Wrap(http.StatusInternalServerError, "Failed to compile checker", err)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return uuid.Nil, apperr.Wrap(http.StatusInternalServerError, "Failed to read compile file", err)
		}
		err = s.s3.UploadFile(ctx, key, bytes.NewReader(body), "application/octet-stream")
		if err != nil {
			return uuid.Nil, apperr.Wrap(http.StatusInternalServerError, "Failed to upload problems to storage", err)
		}
		checkerPath = &key
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
		Content:         content,
		CheckerLanguage: input.CheckerLanguage,
		CheckerPath:     checkerPath,
		TimeLimit:       input.TimeLimit,
		MemoryLimit:     input.MemoryLimit,
		ColorCode:       input.ColorCode,
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
		case pgerrcode.UniqueViolation:
			return uuid.Nil, apperr.Wrap(http.StatusBadRequest, "Duplicate language", err)
		case pgerrcode.InvalidTextRepresentation:
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
