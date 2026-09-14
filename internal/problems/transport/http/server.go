package http

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/f-code-club/rode-battle-api/internal/problems/service"
	"github.com/f-code-club/rode-battle-api/internal/shared"
	"github.com/f-code-club/rode-battle-api/internal/shared/middleware"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Server struct {
	service        service.Service
	accessTokenSvc *shared.TokenService
	s3             *shared.S3Service
	channel        *amqp.Channel
}

func NewServer(
	cfg *shared.Config,
	pool *pgxpool.Pool,
	accessTokenSvc *shared.TokenService,
	s3 *shared.S3Service,
	channel *amqp.Channel,
) Server {
	service := service.New(pool, s3, channel)

	return Server{service, accessTokenSvc, s3, channel}
}

func (s *Server) RegisterRoutes(f *fuego.Server) {
	m := middleware.NewParseToken(s.accessTokenSvc)

	g := fuego.Group(f, "/problems")
	fuego.Post(g, "/", s.CreateProblem,
		option.Middleware(m),
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
