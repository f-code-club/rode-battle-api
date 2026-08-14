package service

import (
	mailservice "github.com/f-code-club/rode-battle-api/internal/mail/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	pool        *pgxpool.Pool
	mailService mailservice.Service
}

func New(pool *pgxpool.Pool, mailService mailservice.Service) Service {
	return Service{pool, mailService}
}
