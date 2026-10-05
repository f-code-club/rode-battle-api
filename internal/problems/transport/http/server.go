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
	service service.Service
	s3      *shared.S3Service
	amqp    *shared.AmqpService
	authSvc auth.Service
}

func NewServer(
	cfg *shared.Config,
	pool *pgxpool.Pool,
	s3 *shared.S3Service,
	amqp *shared.AmqpService,
	judgeUrl string,
	authSvc auth.Service,
) Server {
	service := service.New(pool, s3, amqp, judgeUrl)

	return Server{service, s3, amqp, authSvc}
}

func (s *Server) RegisterRoutes(api huma.API) {
	requireJury := middleware.NewRequireRole(api, s.authSvc, auth.Jury)
	requireAll := middleware.NewRequireRole(api, s.authSvc, auth.Participant, auth.Jury, auth.Admin)

	g := huma.NewGroup(api, "/problems")

	huma.Register(g, huma.Operation{
		OperationID: "problems-create",
		Method:      http.MethodPost,
		Path:        "",
		Summary:     "Create a problem",
		Tags:        []string{"problems"},
		Middlewares: huma.Middlewares{requireJury},
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
		Middlewares: huma.Middlewares{requireAll},
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
		Middlewares: huma.Middlewares{requireAll},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, s.CreateSubmission)

	huma.Register(g, huma.Operation{
		OperationID: "problems-get-all",
		Method:      http.MethodGet,
		Path:        "",
		Summary:     "Get all problems",
		Tags:        []string{"problems"},
		Middlewares: huma.Middlewares{requireJury},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, s.GetProblems)

	huma.Register(g, huma.Operation{
		OperationID: "problem-create-test-case",
		Method:      http.MethodPost,
		Path:        "/{id}/test-case",
		Summary:     "Create new test case for problem",
		Tags:        []string{"problems", "test-case"},
		Middlewares: huma.Middlewares{requireJury},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, s.CreateTestCase)
}
