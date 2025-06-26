package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"wallet-service/internal/domain"
)

// const userId = "a737d022-eabd-4b04-ac0b-87ee9cb10885"

func (s *Server) createWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response(w, http.StatusMethodNotAllowed, ErrHTTPMethod)

		return
	}

	var walletInfo domain.WalletInfo

	ctx := r.Context()
	userInfo := getUserFromContext(ctx)

	if err := json.NewDecoder(r.Body).Decode(&walletInfo); err != nil {
		response(w, http.StatusBadRequest, err.Error())

		return
	}

	wallet := domain.Wallet{
		Id:        uuid.New(),
		UserId:    userInfo.Id.String(),
		Name:      walletInfo.Name,
		Balance:   walletInfo.Balance,
		Currency:  walletInfo.Currency,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	newWallet, err := s.services.CreateWallet(ctx, wallet, userInfo.Id.String())
	if err != nil {
		response(w, http.StatusInternalServerError, err.Error())

		return
	}

	response(w, http.StatusCreated, Map{
		"wallet": newWallet,
	})
}

func (s *Server) getWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response(w, http.StatusMethodNotAllowed, ErrHTTPMethod)

		return
	}

	walletId, err := getWalletId(r)
	if err != nil {
		response(w, http.StatusBadRequest, err.Error())

		return
	}

	ctx := r.Context()
	userInfo := getUserFromContext(ctx)

	wallet, err := s.services.GetWallet(ctx, walletId, userInfo.Id.String())
	if err != nil {
		response(w, http.StatusInternalServerError, err.Error())

		return
	}

	response(w, http.StatusOK, Map{
		"wallet": wallet,
	})
}

func (s *Server) getWallets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response(w, http.StatusMethodNotAllowed, ErrHTTPMethod)

		return
	}

	ctx := r.Context()
	userInfo := getUserFromContext(ctx)

	wallets, err := s.services.GetWallets(ctx, userInfo.Id.String())
	if err != nil {
		response(w, http.StatusInternalServerError, err.Error())

		return
	}

	response(w, http.StatusOK, Map{
		"wallets": wallets,
	})
}

func (s *Server) updateWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		response(w, http.StatusMethodNotAllowed, ErrHTTPMethod)

		return
	}

	walletId, err := getWalletId(r)
	if err != nil {
		response(w, http.StatusBadRequest, err.Error())

		return
	}

	ctx := r.Context()
	userInfo := getUserFromContext(ctx)

	var updateWallet domain.WalletUpdate

	if err := json.NewDecoder(r.Body).Decode(&updateWallet); err != nil {
		response(w, http.StatusBadRequest, err.Error())

		return
	}

	updatedWallet, err := s.services.UpdateWallet(ctx, walletId,
		userInfo.Id.String(), updateWallet)
	if err != nil {
		response(w, http.StatusInternalServerError, err.Error())

		return
	}

	response(w, http.StatusOK, Map{
		"updated wallet": updatedWallet,
	})
}

func (s *Server) deleteWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		response(w, http.StatusMethodNotAllowed, ErrHTTPMethod)

		return
	}

	walletId, err := getWalletId(r)
	if err != nil {
		response(w, http.StatusBadRequest, err.Error())

		return
	}

	ctx := r.Context()
	userInfo := getUserFromContext(ctx)

	if err := s.services.DeleteWallet(ctx, walletId, userInfo.Id.String()); err != nil {
		response(w, http.StatusInternalServerError, err.Error())

		return
	}

	response(w, http.StatusNoContent, nil)
}
