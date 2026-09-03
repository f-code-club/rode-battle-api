package http

import (
	"github.com/f-code-club/rode-battle-api/internal/contests/service"
	"github.com/f-code-club/rode-battle-api/internal/shared"
	"github.com/f-code-club/rode-battle-api/internal/shared/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
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

	return Server{
		service:        service,
		accessTokenSvc: accessTokenSvc,
	}
}

func (s *Server) RegisterRoutes(f *fuego.Server) {
	m := middleware.NewParseToken(s.accessTokenSvc)

	g := fuego.Group(f, "/contests")

	// TODO: RequireRole(JURY)
	fuego.Post(g, "", s.CreateContest,
		option.Middleware(m),
		option.Security(openapi3.SecurityRequirement{"bearerAuth": []string{}}),
	)

	fuego.Get(g, "/{id}/rank", s.GetRank)
	fuego.Get(g, "", s.ListContests)
	fuego.Get(g, "/{id}", s.GetContestDetail)
}
