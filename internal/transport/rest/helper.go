package rest

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// func getUserId(r *http.Request) uuid.UUID {
// 	userId := r.Context().Value("userId").(string)
// 	return uuid.MustParse(userId)
// }

func getWalletId(r *http.Request) (uuid.UUID, error) {
	walletId := mux.Vars(r)["walletId"]

	walletIdParsed, err := uuid.Parse(walletId)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to parse Wallet Id: %w", err)
	}

	return walletIdParsed, nil
}
