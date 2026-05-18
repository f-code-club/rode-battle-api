package auth

import (
	"errors"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	Secret    string `env:"JWT_SECRET" envRequired:"true"`
	ExpiredIn int    `env:"JWT_EXPIRED_IN" envDefault:"5"` // hours
}

type Claims struct {
	jwt.RegisteredClaims
}

var (
	ErrSigningMethod = errors.New("unexpected signing method")
	ErrInvalidClaims = errors.New("invalid token claims")
)

func NewTokenService() (*TokenService, error) {
	t, err := env.ParseAs[TokenService]()
	if err != nil {
		return nil, err
	}

	return &t, nil
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
