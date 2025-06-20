package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"wallet-service/internal/domain"
)

type WalletDB struct {
	db *sqlx.DB
}

func NewWalletRepository(db *sqlx.DB) *WalletDB {
	return &WalletDB{
		db: db,
	}
}

func (w *WalletDB) CreateWallet(ctx context.Context, wallet domain.Wallet, userId string) (domain.Wallet, error) {
	query := `INSERT INTO wallets
	(id, user_id, name, balance, currency, created_at, updated_at, deleted_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	userIdParsed, err := uuid.Parse(userId)
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("failed to parse userId from string to UUID: %w", err)
	}

	_, err = w.db.ExecContext(ctx, query,
		wallet.Id,
		userIdParsed,
		wallet.Name,
		wallet.Balance,
		wallet.Currency,
		wallet.CreatedAt,
		wallet.UpdatedAt,
		wallet.DeletedAt)
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("failed to insert User: %w", err)
	}

	return wallet, nil
}

func (w *WalletDB) GetWallet(ctx context.Context, walletId uuid.UUID, userId string) (domain.Wallet, error) {
	var wallet domain.Wallet

	query := `SELECT id, user_id, name, balance, currency, created_at, updated_at, deleted_at
	FROM wallets
	WHERE id = $1
	AND user_id = $2
	AND deleted_at IS NULL`

	userIdParsed, err := uuid.Parse(userId)
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("failed to parse userId from string to UUID: %w", err)
	}

	if err := w.db.QueryRowContext(ctx, query, walletId, userIdParsed).Scan(
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

func (w *WalletDB) GetWallets(ctx context.Context, userId string) ([]domain.Wallet, error) {
	var wallets []domain.Wallet

	query := `SELECT name, balance, currency, created_at, updated_at, deleted_at
	FROM wallets
	WHERE user_id = $1 AND deleted_at IS NULL`

	userIdParsed, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse userId from string to UUID: %w", err)
	}

	if err := w.db.SelectContext(ctx, &wallets, query, userIdParsed); err != nil {
		return nil, fmt.Errorf("failed to get all wallets: %w", err)
	}

	return wallets, nil
}

func (w *WalletDB) UpdateWallet(ctx context.Context, walletId uuid.UUID, userId string, wallet domain.WalletUpdate) (domain.Wallet, error) {
	newWallet := domain.Wallet{
		Name: wallet.Name,
	}

	userIdParsed, err := uuid.Parse(userId)
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("failed to parse userId from string to UUID: %w", err)
	}

	query := `UPDATE wallets SET name = $1
	WHERE id = $2
	AND user_id = $3
	AND deleted_at IS NULL`

	if _, err := w.db.ExecContext(ctx, query, newWallet.Name, walletId, userIdParsed); err != nil {
		return domain.Wallet{}, fmt.Errorf("failed to update the wallet: %w", err)
	}

	return newWallet, nil
}

func (w *WalletDB) DeleteWallet(ctx context.Context, walletId uuid.UUID, userId string) error {
	query1 := `UPDATE wallets SET deleted_at = NOW()
	WHERE id = $1
	AND user_id = $2
	AND deleted_at IS NULL`

	userIdParsed, err := uuid.Parse(userId)
	if err != nil {
		return fmt.Errorf("failed to parse userId from string to UUID: %w", err)
	}

	if _, err := w.db.ExecContext(ctx, query1, walletId, userIdParsed); err != nil {
		return fmt.Errorf("failed to update deleted_at column: %w", err)
	}

	query2 := `DELETE FROM wallets
	WHERE id = $1 AND user_id = $2`

	if _, err = w.db.ExecContext(ctx, query2, walletId, userIdParsed); err != nil {
		return fmt.Errorf("failed to delete wallet: %w", err)
	}

	return nil
}
