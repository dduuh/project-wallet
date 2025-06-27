package rest

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"time"

	configs "wallet-service/internal/config"
	"wallet-service/internal/repository"
	"wallet-service/internal/service"

	"github.com/gorilla/mux"
)

const maxHeaderBytes = 1 << 20

type Server struct {
	server   *http.Server
	services *service.Service
	userRepo *repository.UsersRepository
	key      *rsa.PublicKey
}

func New(services *service.Service, userRepo *repository.UsersRepository) *Server {
	return &Server{
		services: services,
		userRepo: userRepo,
	}
}

func (s *Server) Run(ctx context.Context, cfg *configs.Config, handler http.Handler) error {
	s.server = &http.Server{
		Addr:           ":" + cfg.HTTP.Port,
		Handler:        handler,
		ReadTimeout:    cfg.HTTP.ReadTimeout,
		WriteTimeout:   cfg.HTTP.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		//nolint:contextcheck
		if err := s.Shutdown(shutdownCtx); err != nil {
			return
		}
	}()

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
	api.Use(s.jwtAuth)

	api.HandleFunc("/wallets", s.getWallets).Methods(http.MethodGet)
	api.HandleFunc("/wallets/{walletId}", s.getWallet).Methods(http.MethodGet)
	api.HandleFunc("/wallets", s.createWallet).Methods(http.MethodPost)
	api.HandleFunc("/wallets/{walletId}", s.updateWallet).Methods(http.MethodPatch)
	api.HandleFunc("/wallets/{walletId}", s.deleteWallet).Methods(http.MethodDelete)

	return r
}
