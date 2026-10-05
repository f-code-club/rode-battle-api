package http

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/f-code-club/rode-battle-api/internal/accounts/service"
	auth "github.com/f-code-club/rode-battle-api/internal/auth/service"
	"github.com/f-code-club/rode-battle-api/internal/shared"
	"github.com/f-code-club/rode-battle-api/internal/shared/middleware"
)

type Server struct {
	service service.Service
	authSvc auth.Service
}

func NewServer(
	cfg *shared.Config,
	pool *pgxpool.Pool,
	authSvc auth.Service,
) Server {
	emailSvc := shared.NewEmailService(cfg.EmailUsername, cfg.EmailPassword, cfg.EmailHost, cfg.EmailPort)
	service := service.New(pool, emailSvc)

	return Server{service, authSvc}
}

func (s *Server) RegisterRoutes(api huma.API) {
	requireAdmin := middleware.NewRequireRole(api, s.authSvc, auth.Admin)

	g := huma.NewGroup(api, "/accounts")

	huma.Register(g, huma.Operation{
		OperationID: "accounts-generate",
		Method:      http.MethodPost,
		Path:        "/generate",
		Summary:     "Generate account",
		Tags:        []string{"accounts"},
		Middlewares: huma.Middlewares{requireAdmin},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, s.Generate)
}
