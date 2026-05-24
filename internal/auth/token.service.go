package auth

import (
	"errors"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrTokenNoSubject          = errors.New("token does not contain subject")
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

func (s *TokenService) ParseToken(tokenStr string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}

		return []byte(s.Secret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}
	if !token.Valid {
		return uuid.Nil, jwt.ErrTokenUnverifiable
	}

	if claims.Subject == "" {
		return uuid.Nil, ErrTokenNoSubject
	}

	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
