package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
	_ "github.com/joho/godotenv/autoload"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	account "github.com/f-code-club/rode-battle-api/internal/accounts/transport/http"
	auth "github.com/f-code-club/rode-battle-api/internal/auth/transport/http"
	contest "github.com/f-code-club/rode-battle-api/internal/contests/transport/http"
	problem "github.com/f-code-club/rode-battle-api/internal/problems/transport/http"
	"github.com/f-code-club/rode-battle-api/internal/shared"
	"github.com/f-code-club/rode-battle-api/internal/shared/middleware"
)

func gracefulShutdown(apiServer *fuego.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	done <- true
}

func build() (*fuego.Server, error) {
	cfg, err := env.ParseAs[shared.Config]()
	if err != nil {
		return nil, err
	}
	pool, err := shared.NewDatabasePool(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	accessTokenSvc := shared.NewTokenService(cfg.JWTAccessSecret, cfg.JWTAccessExpiredIn)

	f := fuego.NewServer(
		fuego.WithAddr(fmt.Sprintf(":%d", cfg.Port)),
		fuego.WithGlobalMiddlewares(middleware.NewCors(cfg.CorsOrigin)),
		fuego.WithEngineOptions(
			fuego.WithErrorHandler(middleware.ValidationErrorHandler),
			fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
				UIHandler: func(specURL string) http.Handler {
					return httpSwagger.Handler(
						httpSwagger.Layout(httpSwagger.StandaloneLayout),
						httpSwagger.PersistAuthorization(true),
						httpSwagger.URL(specURL),
					)
				},
				DisableDefaultServer: true,
				DisableMessages:      true,
				Info: &openapi3.Info{
					Title:       "R.ODE API",
					Description: "R.ODE API",
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
	api := fuego.Group(f, "/api/v1")

	auth := auth.NewServer(&cfg, pool, &accessTokenSvc)
	auth.RegisterRoutes(api)

	account := account.NewServer(&cfg, pool, &accessTokenSvc)
	account.RegisterRoutes(api)

	problem := problem.NewServer(&cfg, pool, &accessTokenSvc)
	problem.RegisterRoutes(api)

	contest := contest.NewServer(pool, &accessTokenSvc)
	contest.RegisterRoutes(f)

	return f, nil
}

func main() {
	server, err := build()
	if err != nil {
		panic(fmt.Sprintf("failed to build server: %s", err))
	}
	done := make(chan bool, 1)

	go gracefulShutdown(server, done)

	err = server.Run()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	<-done
	log.Println("Graceful shutdown complete.")
}
