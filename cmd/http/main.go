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
	"github.com/f-code-club/rode-battle-api/internal/shared"
	"github.com/go-fuego/fuego"
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

	tokenService, err := auth.NewTokenService()
	if err != nil {
		log.Printf("failed to create token service: %v", err)
		return
	}

	authService := auth.NewAuthService(pool, tokenService)

	s, err := server.NewServer(authService)
	if err != nil {
		log.Printf("failed to create new server: %v", err)
		return
	}
	server := s.Build()
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
