package handlers

import (
	"wallet-service/internal/service"

	"github.com/gorilla/mux"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *mux.Router {
	r := mux.NewRouter()

	auth := r.PathPrefix("/auth").Subrouter()
	auth.PathPrefix("/signup").Methods("POST")
	return r
}