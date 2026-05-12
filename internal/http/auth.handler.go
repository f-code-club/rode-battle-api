package http

import (
	"errors"

	"github.com/f-code-club/rode-battle-api/internal/auth"
	"github.com/f-code-club/rode-battle-api/internal/validation"
	"github.com/go-fuego/fuego"
	"github.com/go-playground/validator/v10"
)

func (s *Server) RegisterHandler(
	c fuego.ContextWithBody[RegisterRequest],
) (string, error) {
	body, err := c.Body()
	if err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			mapped := validation.MapValidationError(validationErrs)
			return "", fuego.BadRequestError{
				Detail: mapped.Error(),
			}
		}
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
