package http

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/f-code-club/rode-battle-api/internal/problems/service"
	"github.com/f-code-club/rode-battle-api/internal/shared"
	"github.com/f-code-club/rode-battle-api/internal/shared/middleware"
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
	m := middleware.ParseTokenMiddlewareBuilder{Service: s.accessTokenSvc}.Middleware

	g := fuego.Group(f, "/problems")
	fuego.Get(g, "/{id}", s.GetProblem)
	fuego.Get(g, "/{id}/history", s.GetSubmitHistory)
	fuego.Post(g, "/{id}/submit", s.CreateSubmission,
		option.Middleware(m),
		option.Security(openapi3.SecurityRequirement{"bearerAuth": []string{}}),
	)
}
