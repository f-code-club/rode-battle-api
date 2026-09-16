package service

import (
	"github.com/f-code-club/rode-battle-api/internal/shared"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	pool     *pgxpool.Pool
	s3       *shared.S3Service
	amqp     *shared.AmqpService
	judgeURL string
}

func New(pool *pgxpool.Pool, s3 *shared.S3Service, channel *shared.AmqpService, judgeURL string) Service {
	return Service{pool, s3, channel, judgeURL}
}
