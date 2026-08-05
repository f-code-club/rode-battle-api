package repository

import (
	"context"

	"github.com/f-code-club/rode-battle-api/internal/model"
	"github.com/f-code-club/rode-battle-api/internal/repository/generated"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository struct {
	Pool *pgxpool.Pool
}

func (r *AccountRepository) Create(ctx context.Context, data model.CreateAccountData) (uuid.UUID, error) {
	queries := generated.New(r.Pool)

	return queries.CreateAccount(ctx, generated.CreateAccountParams{
		Email:    data.Email,
		Name:     data.Name,
		Password: data.Password,
		Role:     generated.Role(data.Role),
	})
}

func (r *AccountRepository) GetByEmail(ctx context.Context, email string) (*model.Account, error) {
	queries := generated.New(r.Pool)

	raw, err := queries.GetAccountByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &model.Account{
		ID:       raw.ID,
		Email:    raw.Email,
		Name:     raw.Name,
		Password: raw.Password,
		Role:     model.Role(raw.Role),
		IsBanned: raw.IsBanned,
	}, nil
}
