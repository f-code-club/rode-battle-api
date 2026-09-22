package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/caarlos0/env/v11"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	_ "github.com/joho/godotenv/autoload"

	account "github.com/f-code-club/rode-battle-api/internal/accounts/transport/http"
	authSvc "github.com/f-code-club/rode-battle-api/internal/auth/service"
	auth "github.com/f-code-club/rode-battle-api/internal/auth/transport/http"
	contest "github.com/f-code-club/rode-battle-api/internal/contests/transport/http"
	problem "github.com/f-code-club/rode-battle-api/internal/problems/transport/http"
	"github.com/f-code-club/rode-battle-api/internal/shared"
	"github.com/f-code-club/rode-battle-api/internal/shared/middleware"
)

func gracefulShutdown(server *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	done <- true
}

func build() (*http.Server, error) {
	cfg, err := env.ParseAs[shared.Config]()
	if err != nil {
		return nil, err
	}
	s3Service, err := shared.NewS3Service(context.Background(), shared.S3Config{
		Bucket:   cfg.S3Bucket,
		Region:   cfg.AwsRegion,
		Endpoint: cfg.AwsEndpointURL,
	})
	if err != nil {
		return nil, err
	}
	pool, err := shared.NewDatabasePool(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	amqp, err := shared.NewRabbitService(cfg.AmqpURL, cfg.TaskQueue)
	if err != nil {
		return nil, err
	}
	accessTokenSvc := shared.NewTokenService(cfg.JWTAccessSecret, cfg.JWTAccessExpiredIn)
	refreshTokenSvc := shared.NewTokenService(cfg.JWTRefreshSecret, cfg.JWTRefreshExpiredIn)
	authSvc := authSvc.New(pool, &refreshTokenSvc, &accessTokenSvc)
	out, err := s3Service.Client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	fmt.Println("Connected successfully. Buckets visible to these credentials:")
	for _, b := range out.Buckets {
		fmt.Printf("  - %s\n", aws.ToString(b.Name))
	}

	mux := http.NewServeMux()

	config := huma.DefaultConfig("R.ODE API", "1.0.0")
	config.Info.Description = "R.ODE API"
	config.DocsRenderer = huma.DocsRendererSwaggerUI
	config.DocsRendererConfig = map[string]any{
		"persistAuthorization": true,
	}
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearerAuth": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
			Description:  "Enter your JWT token in the format: Bearer <token>",
		},
	}

	api := humago.New(mux, config)
	apiV1 := huma.NewGroup(api, "/api/v1")

	authServer := auth.NewServer(&cfg, pool, &accessTokenSvc)
	authServer.RegisterRoutes(apiV1)

	accountServer := account.NewServer(&cfg, pool, &accessTokenSvc, authSvc)
	accountServer.RegisterRoutes(apiV1)

	problemServer := problem.NewServer(&cfg, pool, &accessTokenSvc, s3Service, amqp, cfg.JudgeURL, authSvc)
	problemServer.RegisterRoutes(apiV1)

	contestServer := contest.NewServer(pool, &accessTokenSvc, authSvc)
	contestServer.RegisterRoutes(api)

	corsHandler := middleware.NewCors(cfg.CorsOrigin)(mux)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: corsHandler,
	}

	return server, nil
}

func main() {
	server, err := build()
	if err != nil {
		panic(fmt.Sprintf("failed to build server: %s", err))
	}
	done := make(chan bool, 1)

	go gracefulShutdown(server, done)

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	<-done
	log.Println("Graceful shutdown complete.")
}
