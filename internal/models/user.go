package models

import "time"

type User struct {
	ID        uint64    `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password_hash" json:"password"`
	Name      string    `db:"name" json:"name"`
	Status    int       `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UpdateUserRequest struct {
	Name        string `json:"name"`
	UpdateName  bool   `json:"update_name"`
	Email       string `json:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	UserID      uint64 `json:"user_id"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ProfileResponse struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status int    `json:"status"`
}

const (
	// UserStatus
	USER_ACTIVE   = 1
	USER_INACTIVE = 0
)
