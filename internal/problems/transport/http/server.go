package http

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
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

func (s *Server) RegisterRoutes(api huma.API) {
	m := middleware.NewParseToken(api, s.accessTokenSvc)
	requireJury := middleware.NewRequireRole(api, s.authSvc, auth.Jury)

	g := huma.NewGroup(api, "/problems")

	huma.Register(g, huma.Operation{
		OperationID: "problems-create",
		Method:      http.MethodPost,
		Path:        "",
		Summary:     "Create a problem",
		Tags:        []string{"problems"},
		Middlewares: huma.Middlewares{m, requireJury},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, s.CreateProblem)

	huma.Register(g, huma.Operation{
		OperationID: "problems-get",
		Method:      http.MethodGet,
		Path:        "/{id}",
		Summary:     "Get problem detail",
		Tags:        []string{"problems"},
	}, s.GetProblem)

	huma.Register(g, huma.Operation{
		OperationID: "problems-get-history",
		Method:      http.MethodGet,
		Path:        "/{id}/history",
		Summary:     "Get submission history for problem",
		Tags:        []string{"problems"},
		Middlewares: huma.Middlewares{m},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, s.GetSubmitHistory)

	huma.Register(g, huma.Operation{
		OperationID: "problems-create-submission",
		Method:      http.MethodPost,
		Path:        "/{id}/submit",
		Summary:     "Create submission for problem",
		Tags:        []string{"problems"},
		Middlewares: huma.Middlewares{m},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, s.CreateSubmission)
}
