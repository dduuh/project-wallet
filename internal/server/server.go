package server

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	srv *http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.srv = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		WriteTimeout:   time.Minute * 10,
		ReadTimeout:    time.Minute * 10,
		MaxHeaderBytes: 1 << 20,
	}

	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
