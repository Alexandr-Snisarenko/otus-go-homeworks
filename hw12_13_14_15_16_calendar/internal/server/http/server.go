package http

import (
	"context"
)

type Server struct {
	logger Logger
	app    Application
}

func NewServer(logger Logger, app Application) *Server {
	return &Server{
		logger: logger,
		app:    app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	// TODO
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error { //nolint:revive
	// TODO
	return nil
}

// TODO
