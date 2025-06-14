package domain

import "time"

type User struct {
	Id        string     `json:"id"`
	BlockedAt *time.Time `json:"blocked_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
