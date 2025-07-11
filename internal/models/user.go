package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID              string    `json:"id"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	ConfirmPassword string    `json:"confirm_password"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func CreateUserTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(query)
	return err
}

func InsertUser(db *sql.DB, user *User) error {
	query := `INSERT INTO users ( email, password) VALUES (?, ?, ?)`
	_, err := db.Exec(query, user.Email, user.Password)
	return err
}

func LoginUser(db *sql.DB, user *User) (string, error) {
	query := `SELECT * FROM users WHERE email = ? AND password = ?`
	row := db.QueryRow(query, user.Email, user.Password)
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	return user.ID, err
}
