package shared

import (
	"net/http"
	"time"

	"github.com/f-code-club/rode-battle-api/internal/shared/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New(http.StatusUnauthorized, "invalid token")
)

type TokenService struct {
	secret    string
	expiredIn time.Duration
}

func NewTokenService(secret string, expiredIn int) TokenService {
	return TokenService{
		secret:    secret,
		expiredIn: time.Duration(expiredIn) * time.Second,
	}
}

func (s *TokenService) GenerateToken(id uuid.UUID) (string, error) {
	claim := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().
			Add(s.expiredIn)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Subject:  id.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", errors.Wrap(http.StatusInternalServerError, "failed to generate token", err)
	}

	return tokenStr, nil
}

func (s *TokenService) ParseToken(tokenStr string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		return uuid.Nil, errors.Wrap(http.StatusUnauthorized, "invalid token", err)
	}
	if !token.Valid {
		return uuid.Nil, ErrInvalidToken
	}

	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, errors.Wrap(http.StatusUnauthorized, "invalid token", err)
	}

	return id, nil
}
