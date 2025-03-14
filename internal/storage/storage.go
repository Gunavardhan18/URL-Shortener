package storage

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
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
	SaveURL(models.URL) error
	GetURL(shortCode string) (string, error)
	GetAllURLs(userID uint64) ([]*models.URLResponse, error)
	DeleteURL(shortCode string) error
	UpdateURL(shortCode, longURL string) error
	RegisterURLAnalytics(ctx *fiber.Ctx, shortCode string) error
	GetClicks(userId uint64, shortCode string) (uint64, error)
}

func (db *PostgresDB) Ping() error {
	return db.DB.Ping()
}

func (db *PostgresDB) SaveURL(url models.URL) error {
	_, err := db.DB.Exec("INSERT INTO urls (short_code, long_url, user_id, created_at, updated_at, expires_at) VALUES ($1, $2, $3, $4, $5, $6)", url.ShortURL, url.LongURL, url.UserID, time.Now(), time.Now(), url.ExpiresAt)
	return err
}

func (db *PostgresDB) GetURL(shortCode string) (string, error) {
	var longURL string
	err := db.DB.QueryRow("SELECT long_url FROM urls WHERE short_code = $1 and status = $2", shortCode, models.URLActive).Scan(&longURL)
	return longURL, err
}

func (db *PostgresDB) GetAllURLs(userID uint64) ([]*models.URLResponse, error) {
	rows, err := db.DB.Query("SELECT id, short_code, long_url FROM urls WHERE user_id = $1 and status = $2", userID, models.URLActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urls := []*models.URLResponse{}
	for rows.Next() {
		var url models.URLResponse
		err := rows.Scan(&url.ID, &url.ShortURL, &url.LongURL)
		if err != nil {
			return nil, err
		}
		urls = append(urls, &url)
	}
	return urls, nil
}

func (db *PostgresDB) DeleteURL(shortCode string) error {
	_, err := db.DB.Exec("Update urls set status = $1, updated_at = now() WHERE short_code = $2", models.URLInactive, shortCode)
	return err
}

func (db *PostgresDB) UpdateURL(shortCode, longURL string) error {
	_, err := db.DB.Exec("UPDATE urls SET long_url = $1, updated_at = $2 WHERE short_code = $3", longURL, time.Now(), shortCode)
	return err
}

func (db *PostgresDB) RegisterURLAnalytics(ctx *fiber.Ctx, shortCode string) error {
	userID := ctx.Locals("userID").(uint64)
	ip := ctx.IP()
	_, err := db.DB.Exec("INSERT INTO url_analytics (short_code, user_id, ip_address, accessed_at) VALUES ($1, $2, $3, $4)", shortCode, userID, ip, time.Now())
	return err
}

func (db *PostgresDB) GetClicks(userID uint64, shortCode string) (uint64, error) {
	var clicks uint64
	err := db.DB.QueryRow("SELECT count(*) FROM url_analytics WHERE short_code = $1 and user_id = $2", shortCode, userID).Scan(&clicks)
	return clicks, err
}
