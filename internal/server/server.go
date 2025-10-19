package server

import (
	"4imdb-seasons-tracker/internal/config"
	"4imdb-seasons-tracker/internal/handler"
	"context"
	"log"
	"net/http"
)

type Server struct {
	httpServer *http.Server
	logger     *log.Logger
}

func NewServer(cfg config.ServerConfig, handler *handler.Handler, logger *log.Logger) *Server {
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + cfg.Port,
			Handler:      mux,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
		logger: logger,
	}
}

func (s *Server) Start() error {
	s.logger.Printf("Starting server on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
