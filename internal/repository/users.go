package repository

import (
	"context"
	"fmt"
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

func (u *UsersRepository) UpsertUser(ctx context.Context, user domain.User) error {
	query := fmt.Sprintf(`INSERT INTO users
	(id, blocked_at, deleted_at)
	VALUES ($1, $2)
	ON CONFLICT (id) DO UPDATE SET
		blocked_at = excluded.blocked_at,
		deleted_at = excluded.deleted_at`)

	_, err := u.psql.ExecContext(ctx, query, user.BlockedAt, user.DeletedAt)
	if err != nil {
		return err
	}

	return nil
}
