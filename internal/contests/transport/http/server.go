package http

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	auth "github.com/f-code-club/rode-battle-api/internal/auth/service"
	"github.com/f-code-club/rode-battle-api/internal/contests/service"
	"github.com/f-code-club/rode-battle-api/internal/shared"
	"github.com/f-code-club/rode-battle-api/internal/shared/middleware"
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

func (s *Server) RegisterRoutes(api huma.API) {
	m := middleware.NewParseToken(api, s.accessTokenSvc)
	requireJury := middleware.NewRequireRole(api, s.authSvc, auth.Jury)

	g := huma.NewGroup(api, "/contests")

	huma.Register(g, huma.Operation{
		OperationID: "contests-create",
		Method:      http.MethodPost,
		Path:        "",
		Summary:     "Create a new contest",
		Tags:        []string{"contests"},
		Middlewares: huma.Middlewares{m, requireJury},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, s.CreateContest)

	huma.Register(g, huma.Operation{
		OperationID: "contests-get-rank",
		Method:      http.MethodGet,
		Path:        "/{id}/rank",
		Summary:     "Get contest ranking",
		Tags:        []string{"contests"},
	}, s.GetRank)

	huma.Register(g, huma.Operation{
		OperationID: "contests-list",
		Method:      http.MethodGet,
		Path:        "",
		Summary:     "List contests",
		Tags:        []string{"contests"},
	}, s.ListContests)

	huma.Register(g, huma.Operation{
		OperationID: "contests-get-detail",
		Method:      http.MethodGet,
		Path:        "/{id}",
		Summary:     "Get contest detail",
		Tags:        []string{"contests"},
	}, s.GetContestDetail)
}
