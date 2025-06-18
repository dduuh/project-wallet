package domain

import (
	"time"

	"github.com/google/uuid"
)

type WalletId uuid.UUID

type Wallet struct {
	Id        uuid.UUID  `json:"-"         db:"id"`
	UserId    string     `json:"-"         db:"user_id"`
	Name      string     `json:"name"      db:"name"`
	Balance   float64    `json:"balance"   db:"balance"`
	Currency  string     `json:"currency"  db:"currency"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
}

type WalletInfo struct {
	Name     string  `json:"name"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

type WalletUpdate struct {
	Id   string `json:"-"`
	Name string `json:"name"`
}
