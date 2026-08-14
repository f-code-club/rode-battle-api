package service

import (
	mailservice "github.com/f-code-club/rode-battle-api/internal/shared/mailer"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	pool        *pgxpool.Pool
	mailService mailservice.Mailer
}

func New(pool *pgxpool.Pool, mailService mailservice.Mailer) Service {
	return Service{pool, mailService}
}
