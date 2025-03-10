package models

type User struct {
	ID       uint64 `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
	Name     string `db:"name" json:"name"`
	Status   int    `db:"status" json:"status"`
}

type UpdateUserRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type URL struct {
	ID       uint64 `json:"id"`
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
	UserID   uint64 `json:"user_id"`
}

const (
	// UserStatus
	USER_ACTIVE   = 1
	USER_INACTIVE = 0
)
