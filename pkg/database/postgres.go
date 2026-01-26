package database

import (
	"Effective-Mobile-Test/internal/config"
	"database/sql"
	"fmt"
)

func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	conStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Failed to ping database: %v", err)
	}

	return db, nil
}
