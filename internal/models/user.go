package models

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
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
		id TEXT PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(query)
	return err
}

func InsertUser(db *sql.DB, user *User) error {
	query := `INSERT INTO users (id, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(query, user.ID, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	return err
}

func LoginUser(db *sql.DB, user *User) (string, error) {
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE email = ?`
	row := db.QueryRow(query, user.Email)
	var storedPassword string
	err := row.Scan(&user.ID, &user.Email, &storedPassword, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return "", err
	}

	// Compare the provided password with the stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password)); err != nil {
		return "", err
	}

	return user.ID, nil
}
