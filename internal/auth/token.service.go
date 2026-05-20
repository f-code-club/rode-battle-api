package auth

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	Secret    string
	ExpiredIn time.Duration
}

func NewTokenService() (*TokenService, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, err
	}

	return &TokenService{
		Secret:    cfg.Secret,
		ExpiredIn: cfg.ExpiredIn,
	}, nil
}

func (s *TokenService) GenerateToken(userID string) (string, error) {
	claim := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().
			Add(s.ExpiredIn)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Subject:  userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(s.Secret))
}
