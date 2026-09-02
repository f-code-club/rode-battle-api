package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-fuego/fuego"
	"github.com/go-playground/validator/v10"
)

func ValidationErrorHandler(ctx context.Context, err error) error {
	if validateErrs, ok := errors.AsType[validator.ValidationErrors](err); ok {
		errs := make([]fuego.ErrorItem, 0, len(validateErrs))
		for _, err := range validateErrs {
			errs = append(errs, fuego.ErrorItem{
				Name:   err.Tag(),
				Reason: fmt.Sprintf("'%s' violates the '%s' constraint", err.Field(), err.Tag()),
			})
		}

		return fuego.HTTPError{
			Status: http.StatusBadRequest,
			Errors: errs,
		}
	}

	return fuego.ErrorHandler(ctx, err)
}
