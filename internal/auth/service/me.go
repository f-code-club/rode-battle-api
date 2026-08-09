package service

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/f-code-club/rode-battle-api/internal/auth/repository"
	"github.com/f-code-club/rode-battle-api/internal/shared/errors"
)

type Role string

const (
	Participant Role = "participant"
	Jury        Role = "jury"
	Admin       Role = "admin"
)

type Account struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  Role   `json:"role"`
}

func (s *Service) Me(ctx context.Context, id uuid.UUID) (*Account, error) {
	queries := repository.New(s.pool)

	account, err := queries.GetAccount(ctx, id)
	if err != nil {
		return nil, errors.Wrap(http.StatusUnauthorized, "account does not existed", err)
	}

	return &Account{
		Email: account.Email,
		Name:  account.Name,
		Role:  Role(account.Role),
	}, nil
}
