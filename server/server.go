package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/aachex/taskrunner/taskscontroller"
)

type Server struct {
	logger *slog.Logger
	srv    http.Server
}

func New(logger *slog.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

func (s *Server) Start() {
	tasksController := taskscontroller.New(s.logger)
	s.logger.Info("initialized tasksController")

	mux := http.NewServeMux()
	tasksController.RegisterEndpoints(mux)
	s.logger.Info("registered endpoints")

	s.srv = http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
