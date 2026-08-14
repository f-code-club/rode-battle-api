package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/f-code-club/rode-battle-api/internal/account/repository"
	"github.com/f-code-club/rode-battle-api/internal/shared/errors"
)

type Role string

const (
	Participant Role = "participant"
	Jury        Role = "jury"
	Admin       Role = "admin"
)

func (s *Service) Generate(
	ctx context.Context,
	email string,
	name string,
	role Role,
) (uuid.UUID, error) {
	queries := repository.New(s.pool)

	password, err := generatePassword()
	fmt.Printf("Generated password: %s\n", password)

	if err != nil {
		return uuid.Nil, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, errors.Wrap(http.StatusInternalServerError, "failed to create account", err)
	}

	id, err := queries.CreateAccount(ctx, repository.CreateAccountParams{
		Email:    email,
		Name:     name,
		Password: string(hashedPassword),
		Role:     repository.Role(role),
	})
	if err != nil {
		return uuid.Nil, errors.Wrap(http.StatusConflict, "account with given email already existed", err)
	}

	err = s.mailService.SendSingleMail(
		email,
		"Welcome to our platform",
		fmt.Sprintf("Hello %s, your account has been created successfully.", name),
	)
	if err != nil {
		return uuid.Nil, errors.Wrap(http.StatusInternalServerError, "failed to send welcome email", err)
	}

	return id, nil
}

func generatePassword() (string, error) {
	raw, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Wrap(http.StatusInternalServerError, "error generating password", err)
	}
	return raw.String(), nil
}
