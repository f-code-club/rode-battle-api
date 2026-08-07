package shared

import (
	"time"

	"github.com/f-code-club/rode-battle-api/internal/shared/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New(errors.Unauthorized, "invalid token")
)

type TokenService struct {
	Secret    string
	ExpiredIn time.Duration
}

func (s *TokenService) GenerateToken(id uuid.UUID) (string, error) {
	claim := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().
			Add(s.ExpiredIn)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Subject:  id.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString([]byte(s.Secret))
	if err != nil {
		return "", errors.Wrap(errors.Internal, "failed to generate token", err)
	}

	return tokenStr, nil
}

func (s *TokenService) ParseToken(tokenStr string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		return []byte(s.Secret), nil
	})
	if err != nil {
		return uuid.Nil, errors.Wrap(errors.Unauthorized, "invalid token", err)
	}
	if !token.Valid {
		return uuid.Nil, ErrInvalidToken
	}

	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, errors.Wrap(errors.Unauthorized, "invalid token", err)
	}

	return id, nil
}
