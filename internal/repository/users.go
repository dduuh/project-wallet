package repository

import (
	"wallet-service/internal/domain"

	"github.com/jmoiron/sqlx"
)

type UsersRepository struct {
	psql *sqlx.DB
}

func NewUsersRepository(psql *sqlx.DB) *UsersRepository {
	return &UsersRepository{
		psql: psql,
	}
}

func (u *UsersRepository) UpsertUser(user domain.User) (int, error) {
	return 0, nil
}
