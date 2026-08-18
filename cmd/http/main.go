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
	"github.com/go-fuego/fuego"
	_ "github.com/joho/godotenv/autoload"

	account "github.com/f-code-club/rode-battle-api/internal/accounts/transport/http"
	auth "github.com/f-code-club/rode-battle-api/internal/auth/transport/http"
	"github.com/f-code-club/rode-battle-api/internal/shared"
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

	f := shared.NewServer(&cfg)

	auth := auth.NewServer(&cfg, pool, &accessTokenSvc)
	auth.RegisterRoutes(f)

	account := account.NewServer(&cfg, pool, &accessTokenSvc)
	account.RegisterRoutes(f)

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
