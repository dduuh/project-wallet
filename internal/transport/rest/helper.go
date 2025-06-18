package rest

import (
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
	return uuid.Parse(walletId)
}
