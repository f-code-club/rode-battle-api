package service

import (
	"context"
	"net/http"

	"github.com/f-code-club/rode-battle-api/internal/auth/repository"
	"github.com/f-code-club/rode-battle-api/internal/shared/errors"
	"github.com/google/uuid"
)

func (s *Service) HasRole(ctx context.Context, id uuid.UUID, roles ...Role) (bool, error) {
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
