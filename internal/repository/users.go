package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"wallet-service/internal/domain"
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
	query := `INSERT INTO users
	(id, blocked_at, deleted_at)
	VALUES ($1, $2, $3)
	ON CONFLICT (id) DO UPDATE SET
		blocked_at = excluded.blocked_at,
		deleted_at = excluded.deleted_at`

	_, err := u.psql.ExecContext(ctx, query, user.Id, user.BlockedAt, user.DeletedAt)
	if err != nil {
		return fmt.Errorf("failed to UpsertUser: %w", err)
	}

	return nil
}

func (u *UsersRepository) GetUser(ctx context.Context, userId uuid.UUID) (domain.User, error) {
	var user domain.User

	query := `SELECT id, blocked_at, deleted_at FROM users WHERE id = $1`

	if err := u.psql.QueryRowContext(ctx, query, userId).Scan(&user.Id, &user.BlockedAt, &user.DeletedAt); err != nil {
		return domain.User{}, fmt.Errorf("failed to get User: %w", err)
	}

	return user, nil
}
