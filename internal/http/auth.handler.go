package http

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-fuego/fuego"
)

var (
	emailRegex       = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	uppercaseRegex   = regexp.MustCompile(`[A-Z]`)
	lowercaseRegex   = regexp.MustCompile(`[a-z]`)
	numberRegex      = regexp.MustCompile(`\d`)
	specialCharRegex = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)
)

func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	return uppercaseRegex.MatchString(password) &&
		lowercaseRegex.MatchString(password) &&
		numberRegex.MatchString(password) &&
		specialCharRegex.MatchString(password)
}

func validRegisterPayload(body RegisterPayload) error {
	body.Email = strings.TrimSpace(body.Email)

	if body.Email == "" {
		return fmt.Errorf("email is required")
	}

	if !emailRegex.MatchString(body.Email) {
		return fmt.Errorf("invalid email format")
	}

	if body.Password == "" {
		return fmt.Errorf("password is required")
	}

	if !isValidPassword(body.Password) {
		return fmt.Errorf(
			"password must be at least 8 characters and contain uppercase, lowercase, number, and special character",
		)
	}

	return nil
}

func (s *Server) RegisterHandler(c fuego.ContextWithBody[RegisterPayload]) (string, error) {
	body, err := c.Body()
	if err != nil {
		return "", fmt.Errorf("can not parse the body")
	}

	err = validRegisterPayload(body)
	if err != nil {
		return "", fuego.BadRequestError{
			Err:    err,
			Detail: err.Error(),
		}
	}

	id, err := s.auth.Register(c.Context(), body.Email, body.Password, body.Name)
	if err != nil {
		return "", fuego.BadRequestError{
			Err:    err,
			Detail: err.Error(),
		}
	}

	return id, nil
}
