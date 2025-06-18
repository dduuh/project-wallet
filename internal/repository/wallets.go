package repository

import (
	"context"
	"wallet-service/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type walletDB struct {
	db *sqlx.DB
}

func NewWalletRepository(db *sqlx.DB) *walletDB {
	return &walletDB{
		db: db,
	}
}

func (w *walletDB) CreateWallet(ctx context.Context, wallet domain.Wallet, userId string) (domain.Wallet, error) {
	query := `INSERT INTO wallets
	(id, user_id, name, balance, currency, created_at, updated_at, deleted_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	userIdParsed, err := uuid.Parse(userId)
	if err != nil {
		return domain.Wallet{}, err
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
		return domain.Wallet{}, err
	}

	return wallet, nil
}

func (w *walletDB) GetWallet(ctx context.Context, walletId uuid.UUID, userId string) (domain.Wallet, error) {
	var wallet domain.Wallet

	query := `SELECT id, user_id, name, balance, currency, created_at, updated_at, deleted_at
	FROM wallets
	WHERE id = $1
	AND user_id = $2
	AND deleted_at IS NULL`

	userIdParsed, err := uuid.Parse(userId)
	if err != nil {
		return domain.Wallet{}, err
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
		return domain.Wallet{}, err
	}

	return wallet, nil
}

func (w *walletDB) GetWallets(ctx context.Context, userId string) ([]domain.Wallet, error) {
	var wallets []domain.Wallet

	query := `SELECT name, balance, currency, created_at, updated_at, deleted_at
	FROM wallets
	WHERE user_id = $1`

	userIdParsed, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	if err := w.db.SelectContext(ctx, &wallets, query, userIdParsed); err != nil {
		return nil, err
	}

	return wallets, nil
}

func (w *walletDB) UpdateWallet(ctx context.Context, walletId uuid.UUID, userId string, wallet domain.WalletUpdate) (domain.Wallet, error) {
	newWallet := domain.Wallet{
		Name: wallet.Name,
	}

	userIdParsed, err := uuid.Parse(userId)
	if err != nil {
		return domain.Wallet{}, err
	}

	query := `UPDATE wallets SET name = $1
	WHERE id = $2 AND user_id = $3`

	if _, err := w.db.ExecContext(ctx, query, newWallet.Name, walletId, userIdParsed); err != nil {
		return domain.Wallet{}, err
	}

	return newWallet, nil
}

func (w *walletDB) DeleteWallet(ctx context.Context, walletId uuid.UUID, userId string) error {
	query1 := `UPDATE wallets SET deleted_at = NOW()
	WHERE id = $1 AND user_id = $2`

	userIdParsed, err := uuid.Parse(userId)
	if err != nil {
		return err
	}

	if _, err := w.db.ExecContext(ctx, query1, walletId, userIdParsed); err != nil {
		return err
	}

	query2 := `DELETE FROM wallets
	WHERE id = $1 AND user_id = $2`

	if _, err = w.db.ExecContext(ctx, query2, walletId, userIdParsed); err != nil {
		return err
	}

	return nil
}
