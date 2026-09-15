package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/f-code-club/rode-battle-api/internal/problems/repository"
	"github.com/f-code-club/rode-battle-api/internal/shared"
	apperr "github.com/f-code-club/rode-battle-api/internal/shared/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	amqp "github.com/rabbitmq/amqp091-go"
)

type CreateSubmission = repository.CreateSubmissionParams

func (s *Service) CreateSubmission(ctx context.Context, problemID uuid.UUID, accountID uuid.UUID, language Language, code string) (uuid.UUID, error) {
	cfg, err := env.ParseAs[shared.Config]()
	if err != nil {
		return uuid.Nil, apperr.Wrap(http.StatusInternalServerError, "Failed to create submission", err)
	}
	var pgErr *pgconn.PgError
	queries := repository.New(s.pool)

	submissionID, err := queries.CreateSubmission(ctx, CreateSubmission{
		ProblemID: problemID,
		AccountID: accountID,
		Language:  language,
		Code:      code,
	})
	if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
		return uuid.Nil, apperr.Wrap(http.StatusBadRequest, "Invalid language", err)
	}
	if err != nil {
		return uuid.Nil, apperr.Wrap(http.StatusBadRequest, "Failed to create submission", err)
	}
	q, err := s.channel.QueueDeclare(
		cfg.TaskQueue,
		true,
		false,
		false,
		false,
		amqp.Table{
			amqp.QueueTypeArg: amqp.QueueTypeQuorum,
		},
	)
	if err != nil {
		return uuid.Nil, apperr.Wrap(http.StatusInternalServerError, "Failed to declare a queue", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = s.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(submissionID[:]),
	})
	if err != nil {
		return uuid.Nil, apperr.Wrap(http.StatusInternalServerError, "Failed to publish to queue", err)
	}

	return submissionID, nil
}
