package auth

import (
	"context"
	"fmt"

	"github.com/f-code-club/rode-battle-api/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthService struct {
	Pool *pgxpool.Pool
}

func NewAuthService(p *pgxpool.Pool) *AuthService {
	return &AuthService{Pool: p}
}

func (s *AuthService) Register(ctx context.Context, email string, password string, name string) (string, error) {
	queries := database.New(s.Pool)

	_, err := queries.GetAccountByEmail(ctx, email)
	if err == nil {
		return "", fmt.Errorf("The email has been registerd!")
	}

	hashPassword, err := HashPassword(password)
	if err != nil {
		return "", fmt.Errorf("Can not hash the password!")
	}

	id, err := queries.CreateAccount(ctx, database.CreateAccountParams{
		Email:    email,
		Password: hashPassword,
		Name:     name,
	})
	if err != nil {
		return "", fmt.Errorf("Can not create new account!")
	}

	return id.String(), nil
}
