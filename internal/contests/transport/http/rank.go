package http

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	auth "github.com/f-code-club/rode-battle-api/internal/auth/service"
	"github.com/f-code-club/rode-battle-api/internal/contests/service"
	apperr "github.com/f-code-club/rode-battle-api/internal/shared/errors"
	"github.com/f-code-club/rode-battle-api/internal/shared/middleware"
)

type GetRankInput struct {
	ID uuid.UUID `path:"id"`
}

type GetRankOutput struct {
	Body []service.Ranking
}

func (s *Server) GetRank(ctx context.Context, input *GetRankInput) (*GetRankOutput, error) {
	isJury := false
	if accountID, ok := ctx.Value(middleware.AccountIDKey).(uuid.UUID); ok {
		var err error
		isJury, err = s.authSvc.HasRole(ctx, accountID, auth.Jury)
		if err != nil {
			return nil, apperr.Wrap(http.StatusInternalServerError, "Failed to check role", err)
		}
	}

	rankings, err := s.service.GetRank(ctx, input.ID, isJury)
	if err != nil {
		return nil, err
	}
	return &GetRankOutput{Body: rankings}, nil
}
