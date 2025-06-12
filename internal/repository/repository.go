package repository

import (
	"wallet-service/internal/domain"

	"github.com/jmoiron/sqlx"
)

type Auth interface {
	CreateUser(user domain.User) (int, error)
}

type Repository struct {
}

func NewRepository(psql *sqlx.DB) *Repository {
	return &Repository{}
}
