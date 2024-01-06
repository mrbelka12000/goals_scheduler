package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"goals_scheduler/pkg/config"
)

// Connect ..
func Connect(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("./%s?_foreign_keys=on", cfg.PathToDB))
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return db, nil
}
