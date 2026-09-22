package http

import (
	"context"
	"errors"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type CreateContestRequest struct {
	Name     string      `json:"name" required:"true"`
	Start    time.Time   `json:"start" required:"true"`
	End      time.Time   `json:"end" required:"true"`
	Problems []uuid.UUID `json:"problems" required:"true"`
}

type CreateContestInput struct {
	Body CreateContestRequest
}

func (i *CreateContestInput) Resolve(ctx huma.Context) []error {
	var errs []error
	if !i.Body.End.After(i.Body.Start) {
		errs = append(errs, errors.New("end time must be after start time"))
	}
	seen := make(map[uuid.UUID]struct{}, len(i.Body.Problems))
	for _, p := range i.Body.Problems {
		if _, ok := seen[p]; ok {
			errs = append(errs, errors.New("problem IDs must be unique"))
			break
		}
		seen[p] = struct{}{}
	}
	return errs
}

type CreateContestOutput struct {
	Body uuid.UUID
}

func (s *Server) CreateContest(ctx context.Context, input *CreateContestInput) (*CreateContestOutput, error) {
	id, err := s.service.CreateContest(ctx, input.Body.Name, input.Body.Start, input.Body.End, input.Body.Problems)
	if err != nil {
		return nil, err
	}

	return &CreateContestOutput{Body: id}, nil
}
