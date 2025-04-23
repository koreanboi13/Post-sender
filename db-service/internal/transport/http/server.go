package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"db/internal/config"

	"github.com/gorilla/mux"
)

type Server struct {
	server *http.Server
	router *mux.Router
}

func NewServer(config *config.ServerConfig) *Server {
	router := mux.NewRouter()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		server: server,
		router: router,
	}
}

func (s *Server) GetRouter() *mux.Router {
	return s.router
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
