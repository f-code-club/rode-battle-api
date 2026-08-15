package service

import (
	"github.com/f-code-club/rode-battle-api/internal/shared"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	pool     *pgxpool.Pool
	emailSvc shared.EmailService
}

func New(pool *pgxpool.Pool, emailSvc shared.EmailService) Service {
	return Service{pool, emailSvc}
}
