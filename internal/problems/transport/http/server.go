package http

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/jackc/pgx/v5/pgxpool"

	auth "github.com/f-code-club/rode-battle-api/internal/auth/service"
	"github.com/f-code-club/rode-battle-api/internal/problems/service"
	"github.com/f-code-club/rode-battle-api/internal/shared"
	"github.com/f-code-club/rode-battle-api/internal/shared/middleware"
)

type Server struct {
	service        service.Service
	accessTokenSvc *shared.TokenService
	s3             *shared.S3Service
	amqp           *shared.AmqpService
	authSvc        auth.Service
}

func NewServer(
	cfg *shared.Config,
	pool *pgxpool.Pool,
	accessTokenSvc *shared.TokenService,
	s3 *shared.S3Service,
	amqp *shared.AmqpService,
	judgeUrl string,
	authSvc auth.Service,
) Server {
	service := service.New(pool, s3, amqp, judgeUrl)

	return Server{service, accessTokenSvc, s3, amqp, authSvc}
}

func (s *Server) RegisterRoutes(f *fuego.Server) {
	m := middleware.NewParseToken(s.accessTokenSvc)

	requireJury := middleware.NewRequireRole(s.authSvc, auth.Jury)

	g := fuego.Group(f, "/problems")
	fuego.Post(g, "/", s.CreateProblem,
		option.Middleware(m),
		option.Middleware(requireJury),
		option.Security(openapi3.SecurityRequirement{"bearerAuth": []string{}}),
	)
	fuego.Get(g, "/{id}", s.GetProblem)
	fuego.Get(g, "/{id}/history", s.GetSubmitHistory,
		option.Middleware(m),
		option.Security(openapi3.SecurityRequirement{"bearerAuth": []string{}}),
	)
	fuego.Post(g, "/{id}/submit", s.CreateSubmission,
		option.Middleware(m),
		option.Security(openapi3.SecurityRequirement{"bearerAuth": []string{}}),
	)
}
