package service

import (
	"github.com/f-code-club/rode-battle-api/internal/shared"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Service struct {
	pool    *pgxpool.Pool
	s3      *shared.S3Service
	channel *amqp.Channel
}

func New(pool *pgxpool.Pool, s3 *shared.S3Service, channel *amqp.Channel) Service {
	return Service{pool, s3, channel}
}
