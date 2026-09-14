package http

import (
	auth "github.com/f-code-club/rode-battle-api/internal/auth/service"
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
	authSvc        auth.Service
}

func NewServer(
	pool *pgxpool.Pool,
	accessTokenSvc *shared.TokenService,
	authSvc auth.Service,
) Server {
	service := service.New(pool)

	return Server{
		service:        service,
		accessTokenSvc: accessTokenSvc,
		authSvc:        authSvc,
	}
}

func (s *Server) RegisterRoutes(f *fuego.Server) {
	m := middleware.NewParseToken(s.accessTokenSvc)

	requireJury := s.authSvc.RequireRole(auth.Jury)

	g := fuego.Group(f, "/contests")

	fuego.Post(g, "", s.CreateContest,
		option.Middleware(m),
		option.Middleware(requireJury),
		option.Security(openapi3.SecurityRequirement{"bearerAuth": []string{}}),
	)

	fuego.Get(g, "/{id}/rank", s.GetRank)
	fuego.Get(g, "", s.ListContests)
	fuego.Get(g, "/{id}", s.GetContestDetail)
}
