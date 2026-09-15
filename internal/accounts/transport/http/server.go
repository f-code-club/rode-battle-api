package http

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/f-code-club/rode-battle-api/internal/accounts/service"
	auth "github.com/f-code-club/rode-battle-api/internal/auth/service"
	"github.com/f-code-club/rode-battle-api/internal/shared"
	"github.com/f-code-club/rode-battle-api/internal/shared/middleware"
)

type Server struct {
	service        service.Service
	accessTokenSvc *shared.TokenService
	authSvc        auth.Service
}

func NewServer(
	cfg *shared.Config,
	pool *pgxpool.Pool,
	accessTokenSvc *shared.TokenService,
	authSvc auth.Service,
) Server {
	emailSvc := shared.NewEmailService(cfg.EmailUsername, cfg.EmailPassword, cfg.EmailHost, cfg.EmailPort)
	service := service.New(pool, emailSvc)

	return Server{service, accessTokenSvc, authSvc}
}

func (s *Server) RegisterRoutes(f *fuego.Server) {
	m := middleware.NewParseToken(s.accessTokenSvc)

	requireAdmin := middleware.NewRequireRole(s.authSvc, auth.Admin)

	g := fuego.Group(f, "/accounts")
	fuego.Post(g, "/generate", s.Generate,
		option.Middleware(m),
		option.Middleware(requireAdmin),
		option.Security(openapi3.SecurityRequirement{"bearerAuth": []string{}}),
	)
}
