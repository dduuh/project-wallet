package service

import (
	"wallet-service/internal/domain"
	"wallet-service/internal/repository"
)

type Auth interface {
	CreateUser(user domain.User) (int, error)
}

type Service struct {
}

func NewService(repo *repository.Repository) *Service {
	return &Service{}
}
