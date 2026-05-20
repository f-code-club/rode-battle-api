package auth

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	Secret    string
	ExpiredIn int
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
			Add(time.Duration(s.ExpiredIn) * time.Hour)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Subject:  userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(s.Secret))
}
