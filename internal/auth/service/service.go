package service

import (
	"github.com/f-code-club/rode-battle-api/internal/shared"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	pool            *pgxpool.Pool
	refreshTokenSvc *shared.TokenService
	AccessTokenSvc  *shared.TokenService
}

func New(pool *pgxpool.Pool, refreshTokenSvc *shared.TokenService, accessTokenSvc *shared.TokenService) Service {
	return Service{pool, refreshTokenSvc, accessTokenSvc}
}
