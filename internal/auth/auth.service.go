package auth

import (
	"context"
	"errors"

	"github.com/f-code-club/rode-battle-api/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrAccountBanned          = errors.New("account is banned")
	ErrAccountNotVerified     = errors.New("account is not verified")
)

type AuthService struct {
	Pool         *pgxpool.Pool
	TokenService *TokenService
}

type UserBasicProfile struct {
	Email string        `json:"email"`
	Name  string        `json:"name"`
	Role  database.Role `json:"role"`
}

type AuthResponse struct {
	AccessToken string           `json:"access_token"`
	UserProfile UserBasicProfile `json:"user_profile"`
}

func NewAuthService(p *pgxpool.Pool, tokenService *TokenService) *AuthService {
	return &AuthService{
		Pool:         p,
		TokenService: tokenService}
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

func (s *AuthService) Login(ctx context.Context, email string, password string) (*AuthResponse, error) {
	queries := database.New(s.Pool)

	account, err := queries.GetAccountByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if account.IsBanned {
		return nil, ErrAccountBanned
	}

	if !account.IsVerified {
		return nil, ErrAccountNotVerified
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)); err != nil {
		return nil, err
	}

	accessToken, err := s.TokenService.GenerateToken(account.ID.String())
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken: accessToken,
		UserProfile: UserBasicProfile{
			Email: account.Email,
			Name:  account.Name,
			Role:  account.Role,
		},
	}, nil
}

func (s *AuthService) ParseToken(tokenStr string) (uuid.UUID, error) {
	return s.TokenService.ParseToken(tokenStr)
}
