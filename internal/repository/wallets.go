package repository

import (
	"context"
	"errors"
	"fmt"

	"wallet-service/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var ErrUserNotFound = errors.New("user not found")

type WalletDB struct {
	db *sqlx.DB
}

func NewWalletRepository(db *sqlx.DB) *WalletDB {
	return &WalletDB{
		db: db,
	}
}

func (w *WalletDB) CreateWallet(ctx context.Context, wallet domain.Wallet, userId uuid.UUID) (domain.Wallet, error) {
	query := `INSERT INTO wallets
	(id, user_id, name, balance, currency, created_at, updated_at, deleted_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, user_id, name, balance, currency, created_at, updated_at, deleted_at`

	err := w.db.QueryRowContext(ctx, query,
		wallet.Id,
		userId,
		wallet.Name,
		wallet.Balance,
		wallet.Currency,
		wallet.CreatedAt,
		wallet.UpdatedAt,
		wallet.DeletedAt).Scan(
		&wallet.Id,
		&userId,
		&wallet.Name,
		&wallet.Balance,
		&wallet.Currency,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
		&wallet.DeletedAt)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return domain.Wallet{}, ErrUserNotFound
		}

		return domain.Wallet{}, fmt.Errorf("failed to insert wallet: %w", err)
	}

	return wallet, nil
}

func (w *WalletDB) GetWallet(ctx context.Context, walletId uuid.UUID, userId uuid.UUID) (domain.Wallet, error) {
	var wallet domain.Wallet

	query := `SELECT id, user_id, name, balance, currency, created_at, updated_at, deleted_at
	FROM wallets
	WHERE id = $1
	AND user_id = $2
	AND deleted_at IS NULL`

	if err := w.db.QueryRowContext(ctx, query, walletId, userId).Scan(
		&wallet.Id,
		&wallet.UserId,
		&wallet.Name,
		&wallet.Balance,
		&wallet.Currency,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
		&wallet.DeletedAt); err != nil {
		return domain.Wallet{}, fmt.Errorf("failed to get the wallet: %w", err)
	}

	return wallet, nil
}

func (w *WalletDB) GetWallets(ctx context.Context, userId uuid.UUID) ([]domain.Wallet, error) {
	var wallets []domain.Wallet

	query := `SELECT name, balance, currency, created_at, updated_at, deleted_at
	FROM wallets
	WHERE user_id = $1 AND deleted_at IS NULL`

	if err := w.db.SelectContext(ctx, &wallets, query, userId); err != nil {
		return nil, fmt.Errorf("failed to get all wallets: %w", err)
	}

	return wallets, nil
}

func (w *WalletDB) UpdateWallet(ctx context.Context, walletId uuid.UUID, userId uuid.UUID, wallet domain.Wallet) (domain.Wallet, error) {
	newWallet := domain.Wallet{
		Name: wallet.Name,
	}

	query := `UPDATE wallets SET name = $1
	WHERE id = $2
	AND user_id = $3
	AND deleted_at IS NULL`

	if _, err := w.db.ExecContext(ctx, query, newWallet.Name, walletId, userId); err != nil {
		return domain.Wallet{}, fmt.Errorf("failed to update the wallet: %w", err)
	}

	return newWallet, nil
}

func (w *WalletDB) DeleteWallet(ctx context.Context, walletId uuid.UUID, userId uuid.UUID) error {
	tx, err := w.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to create the transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				logrus.Panicf("failed to rollback: %v\n", rErr)
			}
		}
	}()

	query1 := `UPDATE wallets SET deleted_at = NOW()
	WHERE id = $1
	AND user_id = $2
	AND deleted_at IS NULL`

	result, err := tx.ExecContext(ctx, query1, walletId, userId)
	if err != nil {
		return fmt.Errorf("failed to update deleted_at column: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected > 0 {
		query2 := `DELETE FROM wallets
		WHERE id = $1 AND user_id = $2`

		if _, err := tx.ExecContext(ctx, query2, walletId, userId); err != nil {
			return fmt.Errorf("failed to delete wallet: %w", err)
		}
	} else {
		logrus.Warnf("ID wallet (%s) already deleted for the ID user (%s)", walletId, userId)
		return fmt.Errorf("wallet already deleted or not found")
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
