package service

import (
	"github.com/f-code-club/rode-battle-api/internal/shared"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	pool *pgxpool.Pool
	s3   *shared.S3Service
}

func New(pool *pgxpool.Pool, s3 *shared.S3Service) Service {
	return Service{pool, s3}
}
