package tests

import (
	"context"
	"net/http"
	"wallet-service/internal/domain"

	"github.com/google/uuid"
)

const walletPath = "/api/v1/wallets"

var existingUser = domain.User{
	Id: uuid.New(),
}

func (s *IntegrationTestSuite) TestCreateWallet() {
	wallet := domain.Wallet{
		Id:       uuid.New(),
		UserId:   uuid.NewString(),
		Name:     "wallet 1",
		Currency: "USD",
	}

	s.Run("user not found", func() {
		s.sendHTTPRequest(http.MethodGet, walletPath, http.StatusNotFound, &wallet, nil, existingUser)
	})

	s.Run("wallet successfully created", func() {
		err := s.usersRepo.UpsertUser(context.Background(), existingUser)
		s.Require().NoError(err)

		wallet.Id = existingUser.Id

		var createdWallet domain.Wallet

		s.sendHTTPRequest(http.MethodPost, walletPath, http.StatusCreated, &wallet, &createdWallet, existingUser)

		s.Require().Equal(wallet.Id, createdWallet.Id)
		s.Require().Equal(wallet.UserId, createdWallet.UserId)
		s.Require().Equal(wallet.Name, createdWallet.Name)
		s.Require().Equal(0.0, createdWallet.Balance)
		s.Require().Equal(wallet.Currency, createdWallet.Currency)
	})

	s.Run("wallet doesn't belong to the user", func() {
		err := s.usersRepo.UpsertUser(context.Background(), existingUser)
		s.Require().NoError(err)

		otherUser := domain.User{
			Id: uuid.New(),
		}

		err = s.usersRepo.UpsertUser(context.Background(), otherUser)
		s.Require().NoError(err)

		existingUser.Id = otherUser.Id

		s.sendHTTPRequest(http.MethodGet, walletPath, http.StatusNotFound, &wallet, nil, existingUser)
	})
}

func (s *IntegrationTestSuite) TestGetWallet() {
	wallet := domain.Wallet{
		Id: uuid.New(),
		UserId: existingUser.Id.String(),
		Name: "wallet 1",
		Balance: 200.0,
		Currency: "USD",
	}

	err := s.usersRepo.UpsertUser(context.Background(), existingUser)
	s.Require().NoError(err)

	var createdWallet domain.Wallet

	s.sendHTTPRequest(http.MethodPost, walletPath, http.StatusCreated, &wallet, &createdWallet, existingUser)

	s.Run("user not found", func() {		
		s.sendHTTPRequest(http.MethodGet, walletPath, http.StatusNotFound, nil, nil, existingUser)
	})

	s.Run("get wallet successfully", func() {
		walletId := uuid.UUID(createdWallet.Id).String()
		fullWalletPath := walletPath + "/" + walletId

		s.sendHTTPRequest(http.MethodGet, fullWalletPath, http.StatusOK, nil, &createdWallet, existingUser)

		s.Require().Equal(wallet.Id, createdWallet.Id)
		s.Require().Equal(wallet.UserId, createdWallet.UserId)
		s.Require().Equal(wallet.Name, createdWallet.Name)
		s.Require().Equal(wallet.Balance, createdWallet.Balance)
		s.Require().Equal(wallet.Currency, createdWallet.Currency)
	})

	s.Run("wallet not found", func() {
		walletId := uuid.New().String()
		fullWalletPath := walletPath + "/" + walletId

		s.sendHTTPRequest(http.MethodGet, fullWalletPath, http.StatusNotFound, nil, nil, existingUser)
	})

	s.Run("wallet doesn't belong to the user", func() {
		otherUser := domain.User{
			Id: uuid.New(),
		}

		err := s.usersRepo.UpsertUser(context.Background(), otherUser)
		s.Require().NoError(err)

		walletId := uuid.UUID(createdWallet.Id).String()
		fullWalletPath := walletPath + "/" + walletId

		s.sendHTTPRequest(http.MethodGet, fullWalletPath, http.StatusNotFound, nil, nil, otherUser)
	})
}

func (s *IntegrationTestSuite) TestUpdateWallet() {
	wallet := domain.Wallet{
		Id: uuid.New(),
		UserId: existingUser.Id.String(),
		Name: "wallet 1",
		Balance: 250.0,
		Currency: "USD",
	}

	err := s.usersRepo.UpsertUser(context.Background(), existingUser)
	s.Require().NoError(err)

	var createdWallet domain.Wallet

	s.sendHTTPRequest(http.MethodPost, walletPath, http.StatusCreated, &wallet, &createdWallet, existingUser)

	s.Run("user not found", func() {
		nonExistingUser := domain.User{
			Id: uuid.New(),
		}

		s.sendHTTPRequest(http.MethodPatch, walletPath, http.StatusNotFound, nil, nil, nonExistingUser)
	})

	s.Run("wallet not found", func() {
		walletId := uuid.New().String()
		fullWalletPath := walletPath + "/" + walletId

		s.sendHTTPRequest(http.MethodGet, fullWalletPath, http.StatusNotFound, nil, nil, existingUser)
	})

	s.Run("wallet updated successfully", func() {
		updatedWallet := domain.Wallet{
			Name: "Blue frog",
		}

		walletId := uuid.UUID(createdWallet.Id).String()
		fullWalletPath := walletPath + "/" + walletId

		s.sendHTTPRequest(http.MethodPatch, fullWalletPath, http.StatusOK, &updatedWallet, &createdWallet, existingUser)

		s.Require().Equal(updatedWallet.Name, createdWallet.Name)
	})

	s.Run("wallet doesn't belong to the user", func() {
		otherUser := domain.User{
			Id: uuid.New(),
		}

		err := s.usersRepo.UpsertUser(context.Background(), otherUser)
		s.Require().NoError(err)

		walletId := uuid.UUID(createdWallet.Id).String()
		fullWalletPath := walletPath + "/" + walletId

		s.sendHTTPRequest(http.MethodGet, fullWalletPath, http.StatusNotFound, nil, nil, otherUser)
	})
}

func (s *IntegrationTestSuite) TestDeleteWallet() {
	wallet := domain.Wallet{
		Id: uuid.New(),
		UserId: existingUser.Id.String(),
		Name: "wallet 1",
		Balance: 300.0,
		Currency: "USD",
	}

	err := s.usersRepo.UpsertUser(context.Background(), existingUser)
	s.Require().NoError(err)

	var createdWallet domain.Wallet

	s.sendHTTPRequest(http.MethodPost, walletPath, http.StatusCreated, &wallet, &createdWallet, existingUser)

	s.Run("user not found", func() {
		nonExistingUser := domain.User{
			Id: uuid.New(),
		}

		s.sendHTTPRequest(http.MethodGet, walletPath, http.StatusNotFound, nil, nil, nonExistingUser)
	})

	s.Run("wallet not found", func() {
		walletId := uuid.New().String()
		fullWalletPath := walletPath + "/" + walletId

		s.sendHTTPRequest(http.MethodGet, fullWalletPath, http.StatusNotFound, nil, nil, existingUser)
	})

	s.Run("wallet successfully deleted", func() {
		walletId := uuid.UUID(createdWallet.Id).String()
		fullWalletPath := walletPath + "/" + walletId

		s.sendHTTPRequest(http.MethodDelete, fullWalletPath, http.StatusNoContent, nil, nil, existingUser)
	})

	s.Run("wallet doesn't belong to the user", func() {
		otherUser := domain.User{
			Id: uuid.New(),
		}

		err := s.usersRepo.UpsertUser(context.Background(), otherUser)
		s.Require().NoError(err)

		walletId := uuid.UUID(createdWallet.Id).String()
		fullWalletPath := walletPath + "/" + walletId

		s.sendHTTPRequest(http.MethodGet, fullWalletPath, http.StatusNotFound, nil, nil, otherUser)
	})
}

func (s *IntegrationTestSuite) TestGetWallets() {
	err := s.usersRepo.UpsertUser(context.Background(), existingUser)
	s.Require().NoError(err)

	var arrWallets []domain.Wallet

	wallet1 := domain.Wallet{
		Id: uuid.New(),
		UserId: existingUser.Id.String(),
		Name: "wallet 1",
		Balance: 100.0,
		Currency: "USD",
	}
	
	wallet2 := domain.Wallet{
		Id: uuid.New(),
		UserId: existingUser.Id.String(),
		Name: "wallet 2",
		Balance: 50.0,
		Currency: "RUB",
	}
	
	wallet3 := domain.Wallet{
		Id: uuid.New(),
		UserId: existingUser.Id.String(),
		Name: "wallet 3",
		Balance: 1000.0,
		Currency: "USD",
	}
	
	arrWallets = append(arrWallets, wallet1, wallet2, wallet3)

	_, err = s.psql.Database().ExecContext(context.Background(),
		`INSERT INTO wallets (name) VALUES ($1) WHERE id = $2 AND user_id = $3`,
	wallet1.Name, wallet1.Id, wallet1.UserId)
	s.Require().NoError(err)

	_, err = s.psql.Database().ExecContext(context.Background(),
		`INSERT INTO wallets (name) VALUES ($1) WHERE id = $2 AND user_id = $3`,
	wallet2.Name, wallet2.Id, wallet2.UserId)
	s.Require().NoError(err)

	_, err = s.psql.Database().ExecContext(context.Background(),
		`INSERT INTO wallets (name) VALUES ($1) WHERE id = $2 AND user_id = $3`,
	wallet3.Name, wallet3.Id, wallet3.UserId)
	s.Require().NoError(err)

	var created1, created2, created3 domain.Wallet

	s.sendHTTPRequest(http.MethodPost, walletPath, http.StatusCreated, &wallet1, &created1, existingUser)
	s.sendHTTPRequest(http.MethodPost, walletPath, http.StatusCreated, &wallet2, &created2, existingUser)
	s.sendHTTPRequest(http.MethodPost, walletPath, http.StatusCreated, &wallet3, &created3, existingUser)

	s.Run("get successfully all wallets", func() {
		var wallets []domain.Wallet

		s.sendHTTPRequest(http.MethodGet, walletPath, http.StatusOK, nil, &wallets, existingUser)

		s.Require().Len(wallets, len(arrWallets))
	})

	s.Run("user doesn't own any wallets", func() {
		otherUser := domain.User{
			Id: uuid.New(),
		}

		err := s.usersRepo.UpsertUser(context.Background(), otherUser)
		s.Require().NoError(err)

		var wallets []domain.Wallet

		s.sendHTTPRequest(http.MethodGet, walletPath, http.StatusNotFound, nil, &wallets, otherUser)

		s.Require().Len(wallets, 0)
	})
}