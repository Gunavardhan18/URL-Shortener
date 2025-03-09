package models

const (
	// UserRoles
	AdminRole = "admin"
	UserRole  = "user"

	// UserStatus
	ActiveStatus   = "active"
	InactiveStatus = "inactive"

	// URLStatus
	ActiveURLStatus   = "active"
	InactiveURLStatus = "inactive"

	// Error messages
	ErrInvalidCredentials = "invalid credentials"
	ErrInvalidToken       = "invalid token"
	ErrUnauthorized       = "unauthorized"
	ErrInternalServer     = "internal server error"
	ErrInvalidRequest     = "invalid request"
	ErrInvalidURL         = "invalid url"
	ErrURLNotFound        = "url not found"
	ErrUserNotFound       = "user not found"
	ErrUserExists         = "user already exists"
	ErrURLExists          = "url already exists"
	ErrInvalidRole        = "invalid role"
	ErrInvalidStatus      = "invalid status"
	ErrInvalidEmail       = "invalid email"
	ErrInvalidPassword    = "invalid password"
	ErrInvalidName        = "invalid name"
	ErrInvalidShortCode   = "invalid short code"
	ErrInvalidLongURL     = "invalid long url"
	ErrInvalidUserID      = "invalid user id"
	ErrInvalidURLID       = "invalid url id"
	ErrInvalidURLStatus   = "invalid url status"
)
