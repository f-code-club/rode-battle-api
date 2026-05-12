package auth

import (
	"context"
	"errors"
	"time"

	"github.com/f-code-club/rode-battle-api/internal/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var ErrEmailAlreadyRegistered = errors.New("email already registered")

const defaultTimeout = 5 * time.Second

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
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	queries := database.New(s.Pool)

	_, err := queries.GetAccountPasswordByEmail(ctx, email)
	if err != nil {

		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}

	} else {
		return ErrEmailAlreadyRegistered
	}

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
		return err
	}

	return nil
}
