package service

import (
	"wallet-service/internal/repository"
)

type Service struct {
	repo *repository.UsersRepository
}

func NewService(repo *repository.UsersRepository) *Service {
	return &Service{
		repo: repo,
	}
}
