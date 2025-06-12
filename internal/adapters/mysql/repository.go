package mysql

import (
	"database/sql"
	"fmt"
	"sendgridtest/internal/domain"

	_ "github.com/go-sql-driver/mysql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(dsn string) (*Repository, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Repository{db: db}, nil
}

func (r *Repository) SaveEvent(event domain.SendgridEvent) error {

	if event.Email == "" || event.Event == "" || event.Timestamp == 0 {
		return domain.ErrInvalidEvent
	}

	const query = `INSERT INTO sendgrid_events (email, event_type, timestamp) VALUES (?, ?, ?)`
	_, err := r.db.Exec(query, event.Email, event.Event, event.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to save event: %w", err)
	}
	return nil
}

func (r *Repository) GetEventsByEmail(email string) ([]domain.SendgridEvent, error) {

	if email == "" {
		return nil, domain.ErrInvalidEvent
	}

	rows, err := r.db.Query(`
        SELECT email, event_type, timestamp
        FROM sendgrid_events 
        WHERE email = ?
        ORDER BY timestamp DESC
    `, email)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []domain.SendgridEvent
	for rows.Next() {
		var event domain.SendgridEvent
		if err := rows.Scan(&event.Email, &event.Event, &event.Timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return events, nil
}
