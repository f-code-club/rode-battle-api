package http

import (
	"errors"
	"strings"
	"unicode"

	"github.com/f-code-club/rode-battle-api/internal/auth"
	"github.com/go-fuego/fuego"
)

var (
	ErrInvalidPassword = errors.New(
		"password must contain uppercase, lowercase, number, and special character",
	)
	ErrInvalidName   = errors.New("name must contain only letters and spaces")
	ErrInvalidSchool = errors.New("school must contain only letters and spaces")
)

func isValidNameString(s string) bool {
	for _, c := range s {
		if !unicode.IsLetter(c) && c != ' ' {
			return false
		}
	}
	return true
}

func (r *RegisterRequest) isValidPassword() bool {
	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, ch := range r.Password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true

		case unicode.IsLower(ch):
			hasLower = true

		case unicode.IsDigit(ch):
			hasNumber = true

		case strings.ContainsRune(
			`!@#$%^&*()_+-=[]{};':"\|,.<>/?`,
			ch,
		):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

func (r *RegisterRequest) isValidName() bool {
	return isValidNameString(r.Name)
}

func (r *RegisterRequest) isValidSchool() bool {
	return isValidNameString(r.School)
}

func (s *Server) RegisterHandler(
	c fuego.ContextWithBody[RegisterRequest],
) (string, error) {
	body, err := c.Body()
	if err != nil {
		return "", err
	}

	err = s.auth.Register(
		c.Context(),
		body.Email,
		body.Password,
		body.Name,
		body.School,
		body.StudentId,
		body.PhoneNumber,
	)
	if err != nil {
		if errors.Is(err, auth.ErrEmailAlreadyRegistered) {
			return "", fuego.ConflictError{
				Err:    err,
				Detail: err.Error(),
			}
		}
		return "", err
	}

	return "Register successful!", nil
}
