package http

import (
	"github.com/go-fuego/fuego"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/f-code-club/rode-battle-api/internal/problems/service"
	"github.com/f-code-club/rode-battle-api/internal/shared"
)
type Server struct {
	service        service.Service
	accessTokenSvc *shared.TokenService
}

func NewServer(
	cfg *shared.Config,
	pool *pgxpool.Pool,
	accessTokenSvc *shared.TokenService,
) Server {
	service := service.New(pool, accessTokenSvc)

	return Server{service, accessTokenSvc}
}

func (s *Server) RegisterRoutes(f *fuego.Server) {

	g := fuego.Group(f, "/problems")
	fuego.Get(g, "/{id}", s.GetProblem)

}
