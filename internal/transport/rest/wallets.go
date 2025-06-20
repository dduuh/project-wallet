package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"wallet-service/internal/domain"
)

const userId = "a737d022-eabd-4b04-ac0b-87ee9cb10885"

func (h *Server) createWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response(w, http.StatusMethodNotAllowed, ErrHTTPMethod)

		return
	}

	var walletInfo domain.WalletInfo

	userIdConv := uuid.MustParse(userId)

	user, err := h.userRepo.GetUser(r.Context(), userIdConv)
	if err != nil {
		response(w, http.StatusInternalServerError, err.Error())

		return
	}

	if err := json.NewDecoder(r.Body).Decode(&walletInfo); err != nil {
		response(w, http.StatusBadRequest, err.Error())

		return
	}

	wallet := domain.Wallet{
		Id:        uuid.New(),
		UserId:    user.Id.String(),
		Name:      walletInfo.Name,
		Balance:   walletInfo.Balance,
		Currency:  walletInfo.Currency,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	newWallet, err := h.services.CreateWallet(r.Context(), wallet, user.Id.String())
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

	userIdConv := uuid.MustParse(userId)

	user, err := h.userRepo.GetUser(r.Context(), userIdConv)
	if err != nil {
		response(w, http.StatusInternalServerError, err.Error())

		return
	}

	wlt, err := h.services.GetWallet(r.Context(), walletId, user.Id.String())
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

	userIdConv := uuid.MustParse(userId)

	user, err := h.userRepo.GetUser(r.Context(), userIdConv)
	if err != nil {
		response(w, http.StatusInternalServerError, err.Error())

		return
	}

	wallets, err := h.services.GetWallets(r.Context(), user.Id.String())
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

	userIdConv := uuid.MustParse(userId)

	user, err := h.userRepo.GetUser(r.Context(), userIdConv)
	if err != nil {
		response(w, http.StatusInternalServerError, err.Error())

		return
	}

	var updateWallet domain.WalletUpdate

	if err := json.NewDecoder(r.Body).Decode(&updateWallet); err != nil {
		response(w, http.StatusBadRequest, err.Error())

		return
	}

	updatedWallet, err := h.services.UpdateWallet(r.Context(), walletId,
		user.Id.String(), updateWallet)
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

	userIdConv := uuid.MustParse(userId)

	user, err := h.userRepo.GetUser(r.Context(), userIdConv)
	if err != nil {
		response(w, http.StatusInternalServerError, err.Error())

		return
	}

	if err := h.services.DeleteWallet(r.Context(), walletId, user.Id.String()); err != nil {
		response(w, http.StatusInternalServerError, err.Error())

		return
	}

	response(w, http.StatusNoContent, nil)
}
