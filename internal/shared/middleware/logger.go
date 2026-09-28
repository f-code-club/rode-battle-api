package middleware

import (
	"errors"
	"log/slog"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type requestScopeKey struct{}

type requestScope struct {
	err          error
	accountID    uuid.UUID
	hasAccountID bool
}

func scopeFrom(ctx huma.Context) *requestScope {
	scope, _ := ctx.Context().Value(requestScopeKey{}).(*requestScope)
	return scope
}

func NewRequestScope() func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		next(huma.WithValue(ctx, requestScopeKey{}, &requestScope{}))
	}
}

func CaptureErrorTransformer(ctx huma.Context, _ string, v any) (any, error) {
	if v == nil {
		return v, nil
	}
	err, ok := v.(error)
	if !ok {
		return v, nil
	}
	if scope := scopeFrom(ctx); scope != nil {
		scope.err = err
	}
	return v, nil
}

func unwrapError(err error) error {
	for {
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			return err
		}
		err = unwrapped
	}
}

func NewAccessLog(log *slog.Logger) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		start := time.Now()
		next(ctx)

		attrs := []slog.Attr{
			slog.String("method", ctx.Method()),
			slog.String("path", ctx.URL().Path),
			slog.Int("status", ctx.Status()),
			slog.Int64("duration_ms", time.Since(start).Milliseconds()),
		}
		if op := ctx.Operation(); op != nil {
			attrs = append(attrs,
				slog.String("path_template", op.Path),
				slog.String("operation_id", op.OperationID),
			)
		}

		level := slog.LevelInfo
		switch {
		case ctx.Status() >= 500:
			level = slog.LevelError
		case ctx.Status() >= 400:
			level = slog.LevelWarn
		}

		if scope := scopeFrom(ctx); scope != nil {
			if scope.hasAccountID {
				attrs = append(attrs, slog.String("account_id", scope.accountID.String()))
			}
			if scope.err != nil {
				attrs = append(attrs, slog.String("error", scope.err.Error()))
				if cause := unwrapError(scope.err); cause != scope.err && cause.Error() != scope.err.Error() {
					attrs = append(attrs, slog.String("cause", cause.Error()))
				}
			}
		}

		log.LogAttrs(ctx.Context(), level, "request", attrs...)
	}
}
