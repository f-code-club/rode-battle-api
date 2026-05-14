package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/f-code-club/rode-battle-api/internal/auth"
	"github.com/f-code-club/rode-battle-api/internal/database"
	server "github.com/f-code-club/rode-battle-api/internal/http"

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

func main() {
	err := server.Init()
	if err != nil {
		log.Printf("failed to init validator: %v", err)
		return
	}

	pool, err := database.NewPool()
	if err != nil {
		log.Printf("failed to init db pool: %v", err)
		return
	}

	authService := auth.NewAuthService(pool)

	s, err := server.NewServer(authService)
	if err != nil {
		log.Printf("failed to create new server: %v", err)
		return
	}
	server := s.Build()

	done := make(chan bool, 1)

	go gracefulShutdown(server, done)

	err = server.Run()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	<-done
	log.Println("Graceful shutdown complete.")
}
