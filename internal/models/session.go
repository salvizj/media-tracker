package models

import (
	"database/sql"
	"time"
)

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SessionInput struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

func CreateSessionTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT NOT NULL,
		token TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(query)
	return err
}

func InsertSession(db *sql.DB, session *Session) error {
	query := `INSERT INTO sessions (user_id, token) VALUES (?, ?)`
	_, err := db.Exec(query, session.UserID, session.Token)
	return err
}

func GetSessionByToken(db *sql.DB, token string) (*Session, error) {
	query := `SELECT id, user_id, token, created_at, updated_at FROM sessions WHERE token = ?`
	row := db.QueryRow(query, token)

	var s Session
	if err := row.Scan(&s.ID, &s.UserID, &s.Token, &s.CreatedAt, &s.UpdatedAt); err != nil {
		return nil, err
	}
	return &s, nil
}

func DeleteSession(db *sql.DB, id string) error {
	query := `DELETE FROM sessions WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}
func DeleteSessionByToken(db *sql.DB, token string) error {
	query := `DELETE FROM sessions WHERE token = ?`
	_, err := db.Exec(query, token)
	return err
}
