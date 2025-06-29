package domain

import (
	"time"
)

type User struct {
	Id        string     `json:"id" db:"id"`
	BlockedAt *time.Time `json:"blockedAt" db:"blocked_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
}

type UserInfo struct {
	Id string `json:"userId"`
}
