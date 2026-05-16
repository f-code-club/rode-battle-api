package http

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/f-code-club/rode-battle-api/internal/auth"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	Config Config

	auth      *auth.AuthService
	validator *validator.Validate
}

func NewServer(authService *auth.AuthService) (*Server, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, err
	}

	validator, err := NewValidator()
	if err != nil {
		return nil, err
	}

	return &Server{
		Config:    cfg,
		auth:      authService,
		validator: validator,
	}, nil
}

func (s *Server) Build() *fuego.Server {
	f := fuego.NewServer(
		fuego.WithAddr(fmt.Sprintf(":%d", s.Config.Port)),
		fuego.WithGlobalMiddlewares(s.CorsMiddleware),
		fuego.WithValidator(s.validator),
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
