package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"db/internal/config"
	"db/internal/service"
	"db/internal/storage"
	transportHttp "db/internal/transport/http"
)

func main() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	postgresStorage, err := storage.NewPostgres(cfg.Database.GetPostgresConnectionString())
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer postgresStorage.Close()

	chatService := service.NewChatService(postgresStorage)

	server := transportHttp.NewServer(&cfg.Server)
	handler := transportHttp.NewHandler(chatService)
	handler.RegisterRoutes(server.GetRouter())

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server on port %s", cfg.Server.Port)
		if err := server.Start(); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	<-done
	log.Println("Server stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		log.Fatalf("Error stopping server: %v", err)
	}

	log.Println("Server stopped")
}
