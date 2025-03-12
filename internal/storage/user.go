package storage

import (
	"time"

	"github.com/guna/url-shortener/internal/models"
)

type IUserStorage interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(userID uint64) (*models.User, error)
	GetUserByName(name string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User, updatename bool) error
	LogoutUser(userID uint64) error
}

func (db *PostgresDB) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := db.DB.QueryRow("SELECT id, email, password_hash, name, status, created_at, updated_at  FROM users WHERE email = $1 and status = $2", email, models.USER_ACTIVE).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	return &user, err
}

func (db *PostgresDB) GetUserByID(userID uint64) (*models.User, error) {
	var user models.User
	err := db.DB.QueryRow("SELECT id, email, password_hash, name, status, created_at, updated_at FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	return &user, err
}

func (db *PostgresDB) GetUserByName(name string) (*models.User, error) {
	var user models.User
	err := db.DB.QueryRow("SELECT id, email, password_hash, name, status, created_at, updated_at FROM users WHERE name = $1", name).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	return &user, err
}

func (db *PostgresDB) CreateUser(user *models.User) error {
	_, err := db.DB.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)", user.Name, user.Email, user.Password)
	return err
}

func (db *PostgresDB) UpdateUser(user *models.User, updateName bool) error {
	if updateName {
		_, err := db.DB.Exec("UPDATE users SET name = $1, email = $2, password_hash = $3, updated_at = $4 WHERE id = $5", user.Name, user.Email, user.Password, time.Now(), user.ID)
		return err
	}
	_, err := db.DB.Exec("UPDATE users SET email = $1, password_hash = $2, updated_at = $3 WHERE id = $4", user.Email, user.Password, time.Now(), user.ID)
	return err
}

func (db *PostgresDB) LogoutUser(userID uint64) error {
	_, err := db.DB.Exec("UPDATE users set status = $1, updated_at = $2 WHERE id = $3", models.USER_INACTIVE, time.Now(), userID)
	return err
}
