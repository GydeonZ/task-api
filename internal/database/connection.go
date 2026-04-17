package database

import (
	"database/sql"
	"log"

	"github.com/GydeonZ/task-api/internal/config"
	_ "github.com/lib/pq"
)

// DB holds the database connection
var DB *sql.DB

// Connect establishes a connection to the database
func Connect(cfg *config.DatabaseConfig) error {
	var err error
	DB, err = sql.Open("postgres", cfg.DSN)
	if err != nil {
		return err
	}

	// Verify connection
	if err = DB.Ping(); err != nil {
		return err
	}

	// Set connection pool parameters
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	log.Println("Database connected successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
