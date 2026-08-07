package service

import (
	"context"

	"github.com/f-code-club/rode-battle-api/internal/auth/repository"
	"github.com/f-code-club/rode-battle-api/internal/shared"
	"github.com/f-code-club/rode-battle-api/internal/shared/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrAccountBanned = errors.New(errors.Unauthorized, "account with given email is banned")
)

type AuthService struct {
	Pool            *pgxpool.Pool
	RefreshTokenSvc *shared.TokenService
	AccessTokenSvc  *shared.TokenService
}

type TokenPair struct {
	RefreshToken string
	AccessToken  string
}

func (a *AuthService) Login(ctx context.Context, email string, password string) (*TokenPair, error) {
	queries := repository.New(a.Pool)

	account, err := queries.GetAccountByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(errors.Unauthorized, "account with given email does not existed", err)
	}
	if account.IsBanned {
		return nil, ErrAccountBanned
	}
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)); err != nil {
		return nil, errors.Wrap(errors.Unauthorized, "wrong password", err)
	}

	refreshToken, err := a.RefreshTokenSvc.GenerateToken(account.ID)
	if err != nil {
		return nil, err
	}

	accessToken, err := a.AccessTokenSvc.GenerateToken(account.ID)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}

func (a *AuthService) Refresh(ctx context.Context, refreshToken string) (string, error) {
	id, err := a.RefreshTokenSvc.ParseToken(refreshToken)
	if err != nil {
		return "", err
	}

	return a.AccessTokenSvc.GenerateToken(id)
}

type Role string

const (
	Participant Role = "participant"
	Jury        Role = "jury"
	Admin       Role = "admin"
)

type Account struct {
	Email string
	Name  string
	Role  Role
}

func (a *AuthService) Me(ctx context.Context, id uuid.UUID) (*Account, error) {
	queries := repository.New(a.Pool)

	account, err := queries.GetAccount(ctx, id)
	if err != nil {
		return nil, errors.Wrap(errors.Unauthorized, "account does not existed", err)
	}

	return &Account{
		Email: account.Email,
		Name:  account.Name,
		Role:  Role(account.Role),
	}, nil
}
