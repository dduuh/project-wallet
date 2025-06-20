package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	configs "wallet-service/internal/config"
	"wallet-service/internal/repository"
	"wallet-service/internal/service"
)

const maxHeaderBytes = 1 << 20

type Server struct {
	server   *http.Server
	services *service.Service
	userRepo *repository.UsersRepository
}

func New(services *service.Service, userRepo *repository.UsersRepository) *Server {
	return &Server{
		services: services,
		userRepo: userRepo,
	}
}

func (s *Server) Run(cfg *configs.Config, handler http.Handler) error {
	s.server = &http.Server{
		Addr:           ":" + cfg.HTTP.Port,
		Handler:        handler,
		ReadTimeout:    cfg.HTTP.ReadTimeout,
		WriteTimeout:   cfg.HTTP.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	if err := s.server.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to run the HTTP server: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown the HTTP server: %w", err)
	}

	return nil
}

func (s *Server) InitRoutes() *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/wallets", s.getWallets).Methods(http.MethodGet)
	api.HandleFunc("/wallets/{walletId}", s.getWallet).Methods(http.MethodGet)
	api.HandleFunc("/wallets", s.createWallet).Methods(http.MethodPost)
	api.HandleFunc("/wallets/{walletId}", s.updateWallet).Methods(http.MethodPatch)
	api.HandleFunc("/wallets/{walletId}", s.deleteWallet).Methods(http.MethodDelete)

	return r
}
