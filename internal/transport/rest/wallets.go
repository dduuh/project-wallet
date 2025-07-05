package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"wallet-service/internal/domain"
)

var (
	ErrWalletNotFound = errors.New("failed to get the wallet: sql: no rows in result set")
)

func (h *Server) createWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response(w, http.StatusMethodNotAllowed, ErrHTTPMethod)

		return
	}

	var wallet domain.Wallet

	ctx := r.Context()
	userInfo := getUserInfoFromContext(ctx)

	if err := json.NewDecoder(r.Body).Decode(&wallet); err != nil {
		response(w, http.StatusBadRequest, err.Error())

		return
	}

	wallet.CreatedAt = time.Now()
	wallet.UpdatedAt = time.Now()

	user, err := h.userRepo.GetUser(ctx, userInfo.Id)
	if err != nil {
		response(w, http.StatusNotFound, err.Error())

		return
	}

	newWallet, err := h.services.CreateWallet(ctx, wallet, user.Id)
	if err != nil {
		response(w, http.StatusInternalServerError, err.Error())

		return
	}

	response(w, http.StatusCreated, Map{
		"wallet": newWallet,
	})
}

func (h *Server) getWallet(w http.ResponseWriter, r *http.Request) {
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
	userInfo := getUserInfoFromContext(ctx)

	user, err := h.userRepo.GetUser(ctx, userInfo.Id)
	if err != nil {
		response(w, http.StatusNotFound, err.Error())

		return
	}

	wallet, err := h.services.GetWallet(ctx, walletId, user.Id)
	if err != nil {
		if err.Error() == ErrWalletNotFound.Error() {
			response(w, http.StatusNotFound, err.Error())

			return
		}

		response(w, http.StatusInternalServerError, err.Error())

		return
	}

	response(w, http.StatusOK, Map{
		"wallet": wallet,
	})
}

func (h *Server) getWallets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response(w, http.StatusMethodNotAllowed, ErrHTTPMethod)

		return
	}

	ctx := r.Context()
	userInfo := getUserInfoFromContext(ctx)

	user, err := h.userRepo.GetUser(ctx, userInfo.Id)
	if err != nil {
		response(w, http.StatusNotFound, err.Error())

		return
	}

	wallets, err := h.services.GetWallets(ctx, user.Id)
	if err != nil {
		response(w, http.StatusInternalServerError, err.Error())

		return
	}

	response(w, http.StatusOK, Map{
		"wallets": wallets,
	})
}

func (h *Server) updateWallet(w http.ResponseWriter, r *http.Request) {
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
	userInfo := getUserInfoFromContext(ctx)

	user, err := h.userRepo.GetUser(ctx, userInfo.Id)
	if err != nil {
		response(w, http.StatusNotFound, err.Error())

		return
	}

	var updateWallet domain.WalletUpdate

	if err := json.NewDecoder(r.Body).Decode(&updateWallet); err != nil {
		response(w, http.StatusBadRequest, err.Error())

		return
	}

	wlt, err := h.services.GetWallet(ctx, walletId, user.Id)
	if err != nil {
		response(w, http.StatusNotFound, err.Error())

		return
	}

	wholeWallet := domain.Wallet{
		Id:        wlt.Id,
		UserId:    user.Id,
		Name:      updateWallet.Name,
		Currency:  wlt.Currency,
		Balance:   wlt.Balance,
		UpdatedAt: time.Now(),
	}

	updatedWallet, err := h.services.UpdateWallet(ctx, walletId,
		user.Id, wholeWallet)
	if err != nil {
		response(w, http.StatusInternalServerError, err.Error())

		return
	}

	response(w, http.StatusOK, Map{
		"updated wallet": updatedWallet,
	})
}

func (h *Server) deleteWallet(w http.ResponseWriter, r *http.Request) {
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
	userInfo := getUserInfoFromContext(ctx)

	user, err := h.userRepo.GetUser(ctx, userInfo.Id)
	if err != nil {
		response(w, http.StatusNotFound, err.Error())

		return
	}

	if err := h.services.DeleteWallet(ctx, walletId, user.Id); err != nil {
		response(w, http.StatusInternalServerError, err.Error())

		return
	}

	response(w, http.StatusNoContent, nil)
}
