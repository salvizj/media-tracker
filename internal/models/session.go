package models

import (
	"database/sql"
	"time"
)

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func CreateSessionTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		expires_at DATETIME DEFAULT (datetime('now', '+1 day'))
	)`
	_, err := db.Exec(query)
	return err
}

func InsertSession(db *sql.DB, session *Session) error {
	query := `INSERT INTO sessions (id, user_id, created_at, updated_at, expires_at) VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(query, session.ID, session.UserID, session.CreatedAt, session.UpdatedAt, session.ExpiresAt)
	return err
}

func GetSessionBySessionID(db *sql.DB, sessionID string) (*Session, error) {
	query := `SELECT id, user_id, created_at, updated_at, expires_at FROM sessions WHERE id = ? AND expires_at > CURRENT_TIMESTAMP`
	row := db.QueryRow(query, sessionID)

	var s Session
	if err := row.Scan(&s.ID, &s.UserID, &s.CreatedAt, &s.UpdatedAt, &s.ExpiresAt); err != nil {
		return nil, err
	}
	return &s, nil
}

func DeleteSessionBySessionID(db *sql.DB, sessionID string) error {
	query := `DELETE FROM sessions WHERE id = ?`
	_, err := db.Exec(query, sessionID)
	return err
}
func CleanUpExpiredSessions(db *sql.DB) error {
	query := `DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP`
	_, err := db.Exec(query)
	return err
}
