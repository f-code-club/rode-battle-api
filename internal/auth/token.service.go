package auth

import (
	"errors"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	Secret    string
	ExpiredIn int
}

type Claims struct {
	jwt.RegisteredClaims
}

var (
	ErrSigningMethod = errors.New("unexpected signing method")
	ErrInvalidClaims = errors.New("invalid token claims")
)

func NewTokenService() (*TokenService, error) {
	t, err := env.ParseAs[Config]()
	if err != nil {
		return nil, err
	}

	return &TokenService{
		Secret:    t.Secret,
		ExpiredIn: t.ExpiredIn,
	}, nil
}

func (s *TokenService) GenerateToken(userID string) (string, error) {
	claim := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().
				Add(time.Duration(s.ExpiredIn) * time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Subject:  userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(s.Secret))
}
