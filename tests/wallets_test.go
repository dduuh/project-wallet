package tests

import (
	"context"
	"net/http"
	"wallet-service/internal/domain"

	"github.com/google/uuid"
)

const walletPath = `/api/v1/wallets`

var (
	existingUser = domain.User{
		Id: uuid.New(),
	}
)

func (s *IntegrationTestSuite) TestCreateWallet() {
	wallet := domain.Wallet{
		Id:       uuid.New(),
		Name:     "wallet 1",
		Currency: "USD",
	}

	s.Run("user not found", func() {
		s.sendHTTPRequest(http.MethodPost, walletPath, http.StatusNotFound, &wallet, nil, domain.User{
			Id: uuid.New(),
		})
	})

	s.Run("wallet successfully created", func() {
		err := s.usersRepo.UpsertUser(context.Background(), existingUser)
		s.Require().NoError(err)

		wallet.UserId = existingUser.Id

		s.sendHTTPRequest(http.MethodPost, walletPath, http.StatusCreated, &wallet, nil, existingUser)
	})

	s.Run("wallet doesn't belong to the user", func() {
		err := s.usersRepo.UpsertUser(context.Background(), existingUser)
		s.Require().NoError(err)

		otherUser := domain.User{
			Id: uuid.New(),
		}

		otherWallet := domain.Wallet{
			Id:       uuid.New(),
			UserId:   otherUser.Id,
			Name:     wallet.Name,
			Currency: wallet.Currency,
		}

		s.sendHTTPRequest(http.MethodPost, walletPath, http.StatusNotFound, &otherWallet, nil, otherUser)
	})
}

func (s *IntegrationTestSuite) TestGetWallet() {
	wallet := domain.Wallet{
		Id:       uuid.New(),
		Name:     "wallet 1",
		Balance:  200.0,
		Currency: "USD",
	}

	err := s.usersRepo.UpsertUser(context.Background(), existingUser)
	s.Require().NoError(err)

	s.sendHTTPRequest(http.MethodPost, walletPath, http.StatusCreated, &wallet, nil, existingUser)

	wallet.UserId = existingUser.Id

	s.Run("user not found", func() {
		s.sendHTTPRequest(http.MethodHead, walletPath, http.StatusNotFound, &wallet, nil, domain.User{
			Id: uuid.New(),
		})
	})

	s.Run("get wallet successfully", func() {
		fullWalletPath := walletPath + "/" + wallet.Id.String()

		s.sendHTTPRequest(http.MethodGet, fullWalletPath, http.StatusOK, &wallet, nil, existingUser)
	})

	s.Run("wallet not found", func() {
		walletId := uuid.NewString()
		fullWalletPath := walletPath + "/" + walletId

		s.sendHTTPRequest(http.MethodGet, fullWalletPath, http.StatusNotFound, &wallet, nil, existingUser)
	})

	s.Run("wallet doesn't belong to the user", func() {
		otherUser := domain.User{
			Id: uuid.New(),
		}

		err := s.usersRepo.UpsertUser(context.Background(), otherUser)
		s.Require().NoError(err)

		walletId := wallet.Id.String()
		fullWalletPath := walletPath + "/" + walletId

		s.sendHTTPRequest(http.MethodGet, fullWalletPath, http.StatusNotFound, nil, nil, otherUser)
	})
}

func (s *IntegrationTestSuite) TestUpdateWallet() {
	wallet := domain.Wallet{
		Id:       uuid.New(),
		Name:     "wallet 1",
		Balance:  250.0,
		Currency: "USD",
	}

	err := s.usersRepo.UpsertUser(context.Background(), existingUser)
	s.Require().NoError(err)

	s.sendHTTPRequest(http.MethodPost, walletPath, http.StatusCreated, &wallet, nil, existingUser)

	wallet.UserId = existingUser.Id

	s.Run("user not found", func() {
		s.sendHTTPRequest(http.MethodPatch, walletPath, http.StatusNotFound, nil, nil, domain.User{
			Id: uuid.New(),
		})
	})

	s.Run("wallet not found", func() {
		walletId := uuid.New().String()
		fullWalletPath := walletPath + "/" + walletId

		s.sendHTTPRequest(http.MethodGet, fullWalletPath, http.StatusNotFound, nil, nil, existingUser)
	})

	s.Run("wallet updated successfully", func() {
		wallet.Name = "Blue frog"

		walletIdU := wallet.Id.String()
		fullWalletPath := walletPath + "/" + walletIdU

		s.sendHTTPRequest(http.MethodPatch, fullWalletPath, http.StatusOK, &wallet, nil, existingUser)
	})

	s.Run("wallet doesn't belong to the user", func() {
		otherUser := domain.User{
			Id: uuid.New(),
		}

		err := s.usersRepo.UpsertUser(context.Background(), otherUser)
		s.Require().NoError(err)

		walletId := wallet.Id.String()
		fullWalletPath := walletPath + "/" + walletId

		s.sendHTTPRequest(http.MethodGet, fullWalletPath, http.StatusNotFound, nil, nil, otherUser)
	})
}

func (s *IntegrationTestSuite) TestDeleteWallet() {
	wallet := domain.Wallet{
		Id:       uuid.New(),
		UserId:   existingUser.Id,
		Name:     "wallet 1",
		Balance:  300.0,
		Currency: "USD",
	}

	err := s.usersRepo.UpsertUser(context.Background(), existingUser)
	s.Require().NoError(err)

	s.sendHTTPRequest(http.MethodPost, walletPath, http.StatusCreated, &wallet, nil, existingUser)

	wallet.UserId = existingUser.Id

	s.Run("user not found", func() {
		s.sendHTTPRequest(http.MethodDelete, walletPath, http.StatusNotFound, &wallet, nil, existingUser)
	})

	s.Run("wallet not found", func() {
		walletId := uuid.NewString()
		fullWalletPath := walletPath + "/" + walletId

		s.sendHTTPRequest(http.MethodGet, fullWalletPath, http.StatusNotFound, nil, nil, existingUser)
	})

	s.Run("wallet successfully deleted", func() {
		walletId := wallet.Id.String()
		fullWalletPath := walletPath + "/" + walletId

		s.sendHTTPRequest(http.MethodDelete, fullWalletPath, http.StatusNoContent, &wallet, nil, existingUser)
	})

	s.Run("wallet doesn't belong to the user", func() {
		otherUser := domain.User{
			Id: uuid.New(),
		}

		err := s.usersRepo.UpsertUser(context.Background(), otherUser)
		s.Require().NoError(err)

		walletId := wallet.Id.String()
		fullWalletPath := walletPath + "/" + walletId

		s.sendHTTPRequest(http.MethodGet, fullWalletPath, http.StatusNotFound, nil, nil, otherUser)
	})
}

func (s *IntegrationTestSuite) TestGetWallets() {
	err := s.usersRepo.Truncate(context.Background())
	s.Require().NoError(err)

	err = s.usersRepo.UpsertUser(context.Background(), existingUser)
	s.Require().NoError(err)

	wallet1 := domain.Wallet{
		Id:       uuid.New(),
		Name:     "wallet 1",
		Balance:  100.0,
		Currency: "USD",
	}

	wallet2 := domain.Wallet{
		Id:       uuid.New(),
		Name:     "wallet 2",
		Balance:  250.0,
		Currency: "RUB",
	}

	wallet3 := domain.Wallet{
		Id:       uuid.New(),
		Name:     "wallet 3",
		Balance:  1000.0,
		Currency: "USD",
	}

	s.Run("empty list", func() {
		var wallets domain.Wallets

		s.sendHTTPRequest(http.MethodGet, walletPath, http.StatusOK, nil, &wallets, existingUser)

		s.Require().Len(wallets.Wallets, 0)
	})

	wallet1.UserId = existingUser.Id
	created1 := domain.Wallet{}

	s.sendHTTPRequest(http.MethodPost, walletPath, http.StatusCreated, &wallet1, &created1, existingUser)

	wallet2.UserId = existingUser.Id
	created2 := domain.Wallet{}

	s.sendHTTPRequest(http.MethodPost, walletPath, http.StatusCreated, &wallet2, &created2, existingUser)

	wallet3.UserId = existingUser.Id
	created3 := domain.Wallet{}

	s.sendHTTPRequest(http.MethodPost, walletPath, http.StatusCreated, &wallet3, &created3, existingUser)

	s.Run("get successfully all wallets", func() {
		var wallets domain.Wallets

		s.sendHTTPRequest(http.MethodGet, walletPath, http.StatusOK, nil, &wallets, existingUser)

		s.Require().Len(wallets.Wallets, 3)
	})

	s.Run("user doesn't own any wallets", func() {
		otherUser := domain.User{
			Id: uuid.New(),
		}

		var wallets domain.Wallets

		s.sendHTTPRequest(http.MethodGet, walletPath, http.StatusNotFound, nil, nil, otherUser)

		s.Require().Len(wallets.Wallets, 0)
	})
}
