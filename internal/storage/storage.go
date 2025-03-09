package storage

import (
	"database/sql"
	"log"
	"os"

	"github.com/guna/url-shortener/internal/models"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB *sql.DB
}

func NewPostgresDB(DB_URL string) *PostgresDB {
	db, err := sql.Open("postgres", DB_URL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
		os.Exit(1) // Exit immediately
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		os.Exit(1) // Exit immediately
	}

	return &PostgresDB{DB: db}
}

type IStorage interface {
	IURLStorage
	IUserStorage
}

type IURLStorage interface {
	Ping() error
	SaveURL(shortCode, longURL string) error
	GetURL(shortCode string) (string, error)
	GetAllURLs(userID uint64) ([]models.URL, error)
	DeleteURL(shortCode string) error
	UpdateURL(shortCode, longURL string) error
}

func (db *PostgresDB) Ping() error {
	return db.DB.Ping()
}

func (db *PostgresDB) SaveURL(shortCode, longURL string) error {
	_, err := db.DB.Exec("INSERT INTO urls (short_code, long_url) VALUES ($1, $2)", shortCode, longURL)
	return err
}

func (db *PostgresDB) GetURL(shortCode string) (string, error) {
	var longURL string
	err := db.DB.QueryRow("SELECT long_url FROM urls WHERE short_code = $1", shortCode).Scan(&longURL)
	return longURL, err
}

func (db *PostgresDB) GetAllURLs(userID uint64) ([]models.URL, error) {
	rows, err := db.DB.Query("SELECT id, short_code, long_url FROM urls WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urls := []models.URL{}
	for rows.Next() {
		var url models.URL
		err := rows.Scan(&url.ID, &url.ShortURL, &url.LongURL)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func (db *PostgresDB) DeleteURL(shortCode string) error {
	_, err := db.DB.Exec("DELETE FROM urls WHERE short_code = $1", shortCode)
	return err
}

func (db *PostgresDB) UpdateURL(shortCode, longURL string) error {
	_, err := db.DB.Exec("UPDATE urls SET long_url = $1 WHERE short_code = $2", longURL, shortCode)
	return err
}
