package middleware

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"

	apperr "github.com/f-code-club/rode-battle-api/internal/shared/errors"
)

func ServiceErrorTransformer(ctx huma.Context, _ string, v any) (any, error) {
	err, ok := v.(error)
	if !ok {
		return v, nil
	}

	var appErr *apperr.Error
	if !errors.As(err, &appErr) {
		return v, nil
	}

	ctx.SetHeader("Content-Type", "application/problem+json")

	return huma.NewError(appErr.GetStatus(), appErr.DetailMsg()), nil
}
