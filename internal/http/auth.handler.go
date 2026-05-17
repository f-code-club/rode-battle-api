package http

import (
	"errors"

	"github.com/f-code-club/rode-battle-api/internal/auth"
	"github.com/go-fuego/fuego"
)

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
	if errors.Is(err, auth.ErrEmailAlreadyRegistered) {
		return "", fuego.ConflictError{
			Err:    err,
			Detail: err.Error(),
		}
	}
	if err != nil {
		return "", fuego.InternalServerError{
			Err:    err,
			Detail: "internal server error",
		}
	}

	return "Register successful!", nil
}

func (s *Server) LoginHandler(
	c fuego.ContextWithBody[LoginRequest],
) (LoginResponse, error) {
	body, err := c.Body()
	if err != nil {
		return LoginResponse{}, err
	}

	result, err := s.auth.Login(
		c.Context(),
		body.Email,
		body.Password,
	)
	if err != nil {
		return LoginResponse{}, fuego.UnauthorizedError{
			Err:    err,
			Detail: "invalid email or password",
		}
	}

	return LoginResponse{
		AccessToken: result.AccessToken,
		UserProfile: UserBasicProfile{
			Name:  result.Name,
			Email: result.Email,
			Role:  result.Role,
		},
	}, nil
}
