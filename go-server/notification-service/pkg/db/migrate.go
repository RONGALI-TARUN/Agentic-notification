package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return db, nil
}

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			id                     TEXT PRIMARY KEY,
			source                 TEXT NOT NULL,
			event_type             TEXT NOT NULL,
			repository             TEXT,
			attention_class        TEXT NOT NULL DEFAULT 'standard',
			decision_status        TEXT,
			decision_reason        TEXT,
			decision_model_version TEXT,
			decision_confidence    DOUBLE PRECISION,
			decision_matched_rule  TEXT,
			payload                JSONB NOT NULL DEFAULT '{}',
			created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("migrate events table: %w", err)
	}
	return nil
}
