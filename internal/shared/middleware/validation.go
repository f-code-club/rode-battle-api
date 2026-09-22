package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-playground/validator/v10"
)

func ValidationErrorHandler(api huma.API) func(ctx huma.Context, err error) {
	return func(ctx huma.Context, err error) {
		if validateErrs, ok := errors.AsType[validator.ValidationErrors](err); ok {
			errs := make([]error, 0, len(validateErrs))
			for _, valErr := range validateErrs {
				errs = append(errs, fmt.Errorf("'%s' violates the '%s' constraint", valErr.Field(), valErr.Tag()))
			}
			_ = huma.WriteErr(api, ctx, http.StatusBadRequest, "validation failed", errs...)
			return
		}
		_ = huma.WriteErr(api, ctx, http.StatusInternalServerError, "internal server error", err)
	}
}
