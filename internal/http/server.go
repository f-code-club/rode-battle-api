package http

import (
	"fmt"

	"github.com/f-code-club/rode-battle-api/internal/auth"
	"github.com/f-code-club/rode-battle-api/internal/config"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
)

type Server struct {
	Config config.Config

	auth *auth.AuthService
}

func NewServer(
	cfg config.Config,
	authService *auth.AuthService,
) (*Server, error) {
	return &Server{
		Config: cfg,
		auth:   authService,
	}, nil
}

func (s *Server) Build() *fuego.Server {
	f := fuego.NewServer(
		fuego.WithAddr(fmt.Sprintf(":%d", s.Config.Port)),
		fuego.WithGlobalMiddlewares(s.CorsMiddleware),
		fuego.WithEngineOptions(
			fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
				UIHandler:            s.OpenAPIHandler,
				DisableDefaultServer: true,
				DisableMessages:      true,
				Info: &openapi3.Info{
					Title:       "General Service",
					Description: "General Service",
				},
			}),
		),
		fuego.WithSecurity(openapi3.SecuritySchemes{
			"bearerAuth": &openapi3.SecuritySchemeRef{
				Value: openapi3.NewSecurityScheme().
					WithType("http").
					WithScheme("bearer").
					WithBearerFormat("JWT").
					WithDescription("Enter your JWT token in the format: Bearer <token>"),
			},
		}),
	)
	s.RegisterRoutes(f)

	return f
}
