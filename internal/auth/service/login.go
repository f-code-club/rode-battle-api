package service

import (
	"context"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/f-code-club/rode-battle-api/internal/auth/repository"
	"github.com/f-code-club/rode-battle-api/internal/shared/errors"
)

var (
	ErrAccountBanned = errors.New(http.StatusUnauthorized, "account with given email is banned")
)

type TokenPair struct {
	RefreshToken string
	AccessToken  string
}

func (s *Service) Login(ctx context.Context, email string, password string) (*TokenPair, error) {
	queries := repository.New(s.pool)

	account, err := queries.GetAccountByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(http.StatusUnauthorized, "account with given email does not existed", err)
	}
	if account.IsBanned {
		return nil, ErrAccountBanned
	}
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)); err != nil {
		return nil, errors.Wrap(http.StatusUnauthorized, "wrong password", err)
	}

	refreshToken, err := s.refreshTokenSvc.GenerateToken(account.ID)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.AccessTokenSvc.GenerateToken(account.ID)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}
