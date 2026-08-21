package http

import (
	"github.com/f-code-club/rode-battle-api/internal/contest/service"
	"github.com/go-fuego/fuego"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	service service.Service
}

func NewServer(
	pool *pgxpool.Pool,
) Server {
	service := service.New(pool)

	return Server{service}
}

func (s *Server) RegisterRoutes(f *fuego.Server) {

	g := fuego.Group(f, "/contests")

	fuego.Get(g, "/{id}/rank", s.GetRank)
}
