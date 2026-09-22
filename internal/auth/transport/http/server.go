package http

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
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

func (s *Server) RegisterRoutes(api huma.API) {
	m := middleware.NewParseToken(api, s.accessTokenSvc)

	g := huma.NewGroup(api, "/auth")

	huma.Register(g, huma.Operation{
		OperationID: "auth-login",
		Method:      http.MethodPost,
		Path:        "/login",
		Summary:     "Log in with email and password",
		Tags:        []string{"auth"},
	}, s.Login)

	huma.Register(g, huma.Operation{
		OperationID: "auth-refresh",
		Method:      http.MethodGet,
		Path:        "/refresh",
		Summary:     "Refresh access token",
		Tags:        []string{"auth"},
	}, s.Refresh)

	huma.Register(g, huma.Operation{
		OperationID: "auth-logout",
		Method:      http.MethodGet,
		Path:        "/logout",
		Summary:     "Log out",
		Tags:        []string{"auth"},
	}, s.Logout)

	huma.Register(g, huma.Operation{
		OperationID: "auth-me",
		Method:      http.MethodGet,
		Path:        "/me",
		Summary:     "Get current user profile",
		Tags:        []string{"auth"},
		Middlewares: huma.Middlewares{m},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, s.Me)
}
