package service

import (
	"context"
	"fmt"

	"wallet-service/internal/domain"
	"wallet-service/internal/repository"

	"github.com/google/uuid"
)

type wallets interface {
	CreateWallet(ctx context.Context, wallet domain.Wallet, userId uuid.UUID) (domain.Wallet, error)
	GetWallet(ctx context.Context, walletId uuid.UUID, userId uuid.UUID) (domain.Wallet, error)
	GetWallets(ctx context.Context, userId uuid.UUID) ([]domain.Wallet, error)
	UpdateWallet(ctx context.Context, walletId uuid.UUID, userId uuid.UUID, wallet domain.Wallet) (domain.Wallet, error)
	DeleteWallet(ctx context.Context, walletId uuid.UUID, userId uuid.UUID) error
}

type Service struct {
	repo     *repository.UsersRepository
	walletDb wallets
}

func New(repo *repository.UsersRepository, walletDb wallets) *Service {
	return &Service{
		repo:     repo,
		walletDb: walletDb,
	}
}

func (s *Service) CreateWallet(ctx context.Context, wallet domain.Wallet, userId uuid.UUID) (domain.Wallet, error) {
	newWallet, err := s.walletDb.CreateWallet(ctx, wallet, userId)
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("%w", err)
	}

	return newWallet, nil
}

func (s *Service) GetWallet(ctx context.Context, walletId uuid.UUID, userId uuid.UUID) (domain.Wallet, error) {
	wallet, err := s.walletDb.GetWallet(ctx, walletId, userId)
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("%w", err)
	}

	return wallet, nil
}

func (s *Service) GetWallets(ctx context.Context, userId uuid.UUID) ([]domain.Wallet, error) {
	wallets, err := s.walletDb.GetWallets(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return wallets, nil
}

func (s *Service) UpdateWallet(ctx context.Context, walletId uuid.UUID, userId uuid.UUID, wallet domain.Wallet) (domain.Wallet, error) {
	updatedWallet, err := s.walletDb.UpdateWallet(ctx, walletId, userId, wallet)
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("%w", err)
	}

	return updatedWallet, nil
}

func (s *Service) DeleteWallet(ctx context.Context, walletId uuid.UUID, userId uuid.UUID) error {
	err := s.walletDb.DeleteWallet(ctx, walletId, userId)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
