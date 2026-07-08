package event

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type SQLExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

type PostgresRepository struct {
	queryer SQLExecutor
}

func NewPostgresRepository(queryer SQLExecutor) *PostgresRepository {
	return &PostgresRepository{queryer: queryer}
}

func (r *PostgresRepository) Save(evt Event) error {
	payload, err := json.Marshal(evt.Payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	_, err = r.queryer.Exec(
		`INSERT INTO events (id, source, event_type, repository, payload) VALUES ($1, $2, $3, $4, $5)`,
		evt.ID,
		evt.Source,
		evt.EventType,
		evt.Repository,
		string(payload),
	)
	return err
}

func (r *PostgresRepository) List() ([]Event, error) {
	rows, err := r.queryer.Query(`SELECT id, source, event_type, repository, payload FROM events ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Event
	for rows.Next() {
		var evt Event
		var payload string
		if err := rows.Scan(&evt.ID, &evt.Source, &evt.EventType, &evt.Repository, &payload); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(payload), &evt.Payload); err != nil {
			return nil, err
		}
		results = append(results, evt)
	}
	return results, nil
}
