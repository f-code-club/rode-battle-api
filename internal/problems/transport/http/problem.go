package http

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/f-code-club/rode-battle-api/internal/problems/service"
	"github.com/f-code-club/rode-battle-api/internal/shared/middleware"
)

type CreateProblemRequest struct {
	Name            string    `json:"name" required:"true"`
	Content         string    `json:"content" required:"true"`
	CheckerLanguage *Language `json:"checker_language"`
	CheckerCode     *string   `json:"checker_code"`
	TimeLimit       *int32    `json:"time_limit"`
	MemoryLimit     *int32    `json:"memory_limit"`
	Languages       []string  `json:"languages" required:"true" minItems:"1"`
	ColorCode       *string   `json:"color_code"`
}

type CreateProblemInput struct {
	Body CreateProblemRequest
}

var allowedLanguages = map[string]struct{}{
	"rust":   {},
	"cpp":    {},
	"python": {},
	"java":   {},
	"html":   {},
}

func (i *CreateProblemInput) Resolve(ctx huma.Context) []error {
	var errs []error
	seen := make(map[string]struct{}, len(i.Body.Languages))
	for _, l := range i.Body.Languages {
		if _, ok := allowedLanguages[l]; !ok {
			errs = append(errs, errors.New("language must be one of: rust, cpp, python, java, html"))
		}
		if _, ok := seen[l]; ok {
			errs = append(errs, errors.New("languages must be unique"))
			break
		}
		seen[l] = struct{}{}
	}
	return errs
}

type CreateProblemOutput struct {
	Body uuid.UUID
}

type GetProblemInput struct {
	ID uuid.UUID `path:"id"`
}

type GetProblemOutput struct {
	Body *service.Problem
}

func (s *Server) GetProblem(ctx context.Context, input *GetProblemInput) (*GetProblemOutput, error) {
	problem, err := s.service.GetProblem(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return &GetProblemOutput{Body: problem}, nil
}

type GetSubmitHistoryInput struct {
	ID uuid.UUID `path:"id"`
}

type GetSubmitHistoryOutput struct {
	Body []service.ProblemHistory
}

func (s *Server) GetSubmitHistory(ctx context.Context, input *GetSubmitHistoryInput) (*GetSubmitHistoryOutput, error) {
	accountID, ok := ctx.Value(middleware.AccountIDKey).(uuid.UUID)
	if !ok {
		return nil, huma.Error401Unauthorized("unauthorized")
	}

	history, err := s.service.GetSubmitHistory(ctx, input.ID, accountID)
	if err != nil {
		return nil, err
	}

	return &GetSubmitHistoryOutput{Body: history}, nil
}

func (s *Server) CreateProblem(ctx context.Context, input *CreateProblemInput) (*CreateProblemOutput, error) {
	id, err := s.service.CreateProblem(ctx, service.CreateProblemInput{
		Name:            input.Body.Name,
		Content:         input.Body.Content,
		CheckerLanguage: input.Body.CheckerLanguage,
		CheckerCode:     input.Body.CheckerCode,
		TimeLimit:       input.Body.TimeLimit,
		MemoryLimit:     input.Body.MemoryLimit,
		ColorCode:       input.Body.ColorCode,
	}, input.Body.Languages)
	if err != nil {
		return nil, err
	}

	return &CreateProblemOutput{Body: id}, nil
}
