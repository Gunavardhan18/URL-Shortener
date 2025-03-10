package storage

import "github.com/guna/url-shortener/internal/models"

type IUserStorage interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(userID uint64) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	LogoutUser(userID uint64) error
}

func (db *PostgresDB) GetUserByEmail(email string) (*models.User, error) {
	var user *models.User
	err := db.DB.QueryRow("SELECT * FROM users WHERE email = $1 and status = $2", email).Scan(&user)
	return user, err
}

func (db *PostgresDB) GetUserByID(userID uint64) (*models.User, error) {
	var user *models.User
	err := db.DB.QueryRow("SELECT * FROM users WHERE id = $1 and status = $2", userID).Scan(&user)
	return user, err
}

func (db *PostgresDB) CreateUser(user *models.User) error {
	_, err := db.DB.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, user.Password)
	return err
}

func (db *PostgresDB) UpdateUser(user *models.User) error {
	_, err := db.DB.Exec("UPDATE users SET email = $1, password = $2 WHERE id = $3", user.Email, user.Password, user.ID)
	return err
}

func (db *PostgresDB) LogoutUser(userID uint64) error {
	_, err := db.DB.Exec("UPDATE users set status = $1 WHERE id = $2", models.USER_INACTIVE, userID)
	return err
}
