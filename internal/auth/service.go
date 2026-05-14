package auth

import (
	"context"
	"errors"

	"github.com/f-code-club/rode-battle-api/internal/database"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var ErrEmailAlreadyRegistered = errors.New("email already registered")

type AuthService struct {
	Pool *pgxpool.Pool
}

func NewAuthService(p *pgxpool.Pool) *AuthService {
	return &AuthService{Pool: p}
}

func (s *AuthService) Register(ctx context.Context,
	email string, password string, name string,
	school string, studentID string, phoneNumber string,
) error {
	queries := database.New(s.Pool)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = queries.CreateAccount(ctx, database.CreateAccountParams{
		Email:       email,
		Password:    string(hashedPassword),
		Name:        name,
		School:      &school,
		StudentID:   &studentID,
		PhoneNumber: &phoneNumber,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return ErrEmailAlreadyRegistered
		}
		return err
	}
	return nil
}
