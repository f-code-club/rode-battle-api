package service

import (
	"context"
	"net/http"

	"github.com/f-code-club/rode-battle-api/internal/auth/repository"
	"github.com/f-code-club/rode-battle-api/internal/shared/errors"
	"github.com/f-code-club/rode-battle-api/internal/shared/middleware"
	"github.com/google/uuid"
)

func (s *Service) RequireRole(roles ...Role) func(http.Handler) http.Handler {
	return middleware.NewRequireRole(func(ctx context.Context, id uuid.UUID) (bool, error) {
		return s.hasRole(ctx, id, roles...)
	})
}

func (s *Service) hasRole(ctx context.Context, id uuid.UUID, roles ...Role) (bool, error) {
	queries := repository.New(s.pool)

	role, err := queries.GetRoleById(ctx, id)
	if err != nil {
		return false, errors.Wrap(http.StatusUnauthorized, "account does not existed", err)
	}

	userRole := Role(role)
	for _, requiredRole := range roles {
		if userRole == requiredRole {
			return true, nil
		}
	}

	return false, nil
}
