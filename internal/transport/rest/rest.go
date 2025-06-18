package rest

import (
	"context"
	"fmt"
	"net/http"
	configs "wallet-service/internal/config"
	"wallet-service/internal/repository"
	"wallet-service/internal/service"

	"github.com/gorilla/mux"
)

type Handler struct {
	server   *http.Server
	services *service.Service
	userRepo *repository.UsersRepository
}

func New(services *service.Service, userRepo *repository.UsersRepository) *Handler {
	return &Handler{
		services: services,
		userRepo: userRepo,
	}
}

func (h *Handler) Run(cfg *configs.Config, handler http.Handler) error {
	h.server = &http.Server{
		Addr:           fmt.Sprintf(":%s", cfg.HTTPCfg.Port),
		Handler:        handler,
		ReadTimeout:    cfg.HTTPCfg.ReadTimeout,
		WriteTimeout:   cfg.HTTPCfg.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	return h.server.ListenAndServe()
}

func (h *Handler) Shutdown(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}

func (h *Handler) InitRoutes() *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/wallets", h.getWallets).Methods("GET")
	api.HandleFunc("/wallets/{walletId}", h.getWallet).Methods("GET")
	api.HandleFunc("/wallets", h.createWallet).Methods("POST")
	api.HandleFunc("/wallets/{walletId}", h.updateWallet).Methods("PATCH")
	api.HandleFunc("/wallets/{walletId}", h.deleteWallet).Methods("DELETE")

	return r
}
