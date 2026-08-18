package http

import (
	"github.com/f-code-club/rode-battle-api/internal/contest/service"
	"github.com/f-code-club/rode-battle-api/internal/shared"
	"github.com/go-fuego/fuego"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	service        service.Service
	accessTokenSvc *shared.TokenService
}

func NewServer(
	pool *pgxpool.Pool,
	accessTokenSvc *shared.TokenService,
) Server {
	service := service.New(pool)

	return Server{service, accessTokenSvc}
}

func (s *Server) RegisterRoutes(f *fuego.Server) {

	g := fuego.Group(f, "/contest")

	fuego.Get(g, "/{id}/rank", s.Rank)
}
