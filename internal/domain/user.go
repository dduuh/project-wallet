package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID  `json:"id" db:"id"`
	BlockedAt *time.Time `json:"blockedAt" db:"blocked_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
}

type UserInfo struct {
	Id uuid.UUID `json:"userId"`
}
