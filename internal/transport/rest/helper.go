package rest

import (
	"fmt"
	"net/http"
	"context"
	"wallet-service/internal/domain"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func getWalletId(r *http.Request) (uuid.UUID, error) {
	walletId := mux.Vars(r)["walletId"]

	walletIdParsed, err := uuid.Parse(walletId)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to parse Wallet Id: %w", err)
	}

	return walletIdParsed, nil
}

func getUserFromContext(ctx context.Context) domain.UserInfo {
	userInfo := ctx.Value(userIdContext).(domain.UserInfo)
	
	return userInfo
}