package rest

import (
	"context"
	"crypto/rsa"
	"errors"
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

func New(services *service.Service, userRepo *repository.UsersRepository, key *rsa.PublicKey) *Server {
	return &Server{
		services: services,
		userRepo: userRepo,
		key:      key,
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

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to run the HTTP server: %w", err)
	}

<<<<<<< HEAD
	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		//nolint:contextcheck
		if err := s.Shutdown(shutdownCtx); err != nil {
			return
		}
	}()

=======
	shutdownCtx, cancel := context.WithTimeout(ctx, 5 * time.Second)
	defer cancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("failed to shutdown the HTTP server: %w", err)
	}
	
>>>>>>> 75e84a18a119dcf4c0fc171fbe504bdb132894dd
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
