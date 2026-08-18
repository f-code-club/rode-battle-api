package http

import (
	"github.com/go-fuego/fuego"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/f-code-club/rode-battle-api/internal/accounts/service"
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
	emailSvc := shared.NewEmailService(cfg.EmailUsername, cfg.EmailPassword, cfg.EmailHost, cfg.EmailPort)
	service := service.New(pool, emailSvc)

	return Server{service, accessTokenSvc}
}

func (s *Server) RegisterRoutes(f *fuego.Server) {
	// m := middleware.ParseTokenMiddlewareBuilder{Service: s.accessTokenSvc}.Middleware

	g := fuego.Group(f, "/accounts")
	fuego.Post(g, "/generate", s.Generate) //TODO: add authorization
}
