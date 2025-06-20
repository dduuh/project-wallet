package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"wallet-service/internal/domain"
	"wallet-service/internal/repository"
)

var (
	ErrCreateWallet = errors.New("failed to create the wallet")
	ErrGetWallet    = errors.New("failed to get the wallet")
	ErrGetWallets   = errors.New("failed to get a wallets")
	ErrUpdateWallet = errors.New("failed to update the wallet")
	ErrDeleteWallet = errors.New("failed to delete the wallet")
)

type wallets interface {
	CreateWallet(ctx context.Context, wallet domain.Wallet, userId string) (domain.Wallet, error)
	GetWallet(ctx context.Context, walletId uuid.UUID, userId string) (domain.Wallet, error)
	GetWallets(ctx context.Context, userId string) ([]domain.Wallet, error)
	UpdateWallet(ctx context.Context, walletId uuid.UUID, userId string, wallet domain.WalletUpdate) (domain.Wallet, error)
	DeleteWallet(ctx context.Context, walletId uuid.UUID, userId string) error
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

func (s *Service) CreateWallet(ctx context.Context, wallet domain.Wallet, userId string) (domain.Wallet, error) {
	newWallet, err := s.walletDb.CreateWallet(ctx, wallet, userId)
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("%w: %w", ErrCreateWallet, err)
	}

	return newWallet, nil
}

func (s *Service) GetWallet(ctx context.Context, walletId uuid.UUID, userId string) (domain.Wallet, error) {
	wallet, err := s.walletDb.GetWallet(ctx, walletId, userId)
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("%w: %w", ErrGetWallet, err)
	}

	return wallet, nil
}

func (s *Service) GetWallets(ctx context.Context, userId string) ([]domain.Wallet, error) {
	wallets, err := s.walletDb.GetWallets(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGetWallets, err)
	}

	return wallets, nil
}

func (s *Service) UpdateWallet(ctx context.Context, walletId uuid.UUID, userId string, wallet domain.WalletUpdate) (domain.Wallet, error) {
	updatedWallet, err := s.walletDb.UpdateWallet(ctx, walletId, userId, wallet)
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("%w: %w", ErrUpdateWallet, err)
	}

	return updatedWallet, nil
}

func (s *Service) DeleteWallet(ctx context.Context, walletId uuid.UUID, userId string) error {
	err := s.walletDb.DeleteWallet(ctx, walletId, userId)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDeleteWallet, err)
	}

	return nil
}
