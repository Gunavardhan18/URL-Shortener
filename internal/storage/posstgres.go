package storage

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB *sql.DB
}

func NewPostgresDB(DB_URL string) *PostgresDB {
	db, err := sql.Open("postgres", DB_URL)
	if err != nil {
		log.Fatal(err)
	}
	return &PostgresDB{DB: db}
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
