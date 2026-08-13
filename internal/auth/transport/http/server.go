package http

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/f-code-club/rode-battle-api/internal/auth/service"
	"github.com/f-code-club/rode-battle-api/internal/shared"
	"github.com/f-code-club/rode-battle-api/internal/shared/middleware"
)

const refreshTokenCookie = "refresh_token"

type Server struct {
	service        service.Service
	accessTokenSvc *shared.TokenService
}

func NewServer(
	cfg *shared.Config,
	pool *pgxpool.Pool,
	accessTokenSvc *shared.TokenService,
) Server {
	refreshTokenSvc := shared.NewTokenService(cfg.JWTRefreshSecret, cfg.JWTRefreshExpiredIn)
	service := service.New(pool, &refreshTokenSvc, accessTokenSvc)

	return Server{service, accessTokenSvc}
}

func (s *Server) RegisterRoutes(f *fuego.Server) {
	m := middleware.ParseTokenMiddlewareBuilder{Service: s.accessTokenSvc}.Middleware

	g := fuego.Group(f, "/auth")
	fuego.Post(g, "/login", s.Login)
	fuego.Get(g, "/refresh", s.Refresh)
	fuego.Get(g, "/me", s.Me,
		option.Middleware(m),
		option.Security(openapi3.SecurityRequirement{"bearerAuth": []string{}}),
	)
}
