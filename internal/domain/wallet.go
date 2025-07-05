package domain

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	Id        uuid.UUID  `json:"id"        db:"id"`
	UserId    uuid.UUID  `json:"userId"    db:"user_id"`
	Name      string     `json:"name"      db:"name"`
	Balance   float64    `json:"balance"   db:"balance"`
	Currency  string     `json:"currency"  db:"currency"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
}

type Wallets struct {
	Wallets []Wallet `json:"wallets"`
}

type WalletUpdate struct {
	Name string `json:"name"`
}
