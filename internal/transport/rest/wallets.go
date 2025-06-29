package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"wallet-service/internal/domain"
)

func (h *Server) createWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response(w, http.StatusMethodNotAllowed, ErrHTTPMethod)

		return
	}

	var walletInfo domain.WalletInfo

	ctx := r.Context()
	userInfo := getUserInfoFromContext(ctx)

	if err := json.NewDecoder(r.Body).Decode(&walletInfo); err != nil {
		response(w, http.StatusBadRequest, err.Error())

		return
	}

	wallet := domain.Wallet{
		Id:        uuid.New(),
		UserId:    userInfo.Id,
		Name:      walletInfo.Name,
		Balance:   walletInfo.Balance,
		Currency:  walletInfo.Currency,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	newWallet, err := h.services.CreateWallet(ctx, wallet, userInfo.Id)
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

	wlt, err := h.services.GetWallet(ctx, walletId, userInfo.Id)
	if err != nil {
		response(w, http.StatusInternalServerError, err.Error())

		return
	}

	response(w, http.StatusOK, Map{
		"wallet": wlt,
	})
}

func (h *Server) getWallets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response(w, http.StatusMethodNotAllowed, ErrHTTPMethod)

		return
	}

	ctx := r.Context()
	userInfo := getUserInfoFromContext(ctx)

	wallets, err := h.services.GetWallets(ctx, userInfo.Id)
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

	var updateWallet domain.WalletUpdate

	if err := json.NewDecoder(r.Body).Decode(&updateWallet); err != nil {
		response(w, http.StatusBadRequest, err.Error())

		return
	}

	updatedWallet, err := h.services.UpdateWallet(ctx, walletId,
		userInfo.Id, updateWallet)
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

	if err := h.services.DeleteWallet(ctx, walletId, userInfo.Id); err != nil {
		response(w, http.StatusInternalServerError, err.Error())

		return
	}

	response(w, http.StatusNoContent, nil)
}
