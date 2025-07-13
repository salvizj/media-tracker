package models

import (
	"database/sql"
	"time"
)

type Movie struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Date      string    `json:"date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    string    `json:"user_id"`
}

type MovieInput struct {
	Name   string `json:"name"`
	UserID string `json:"user_id"`
}

func CreateMoviesTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS movies (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		date TEXT,
		user_id TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(query)
	return err
}

func InsertMovie(db *sql.DB, m *Movie) error {
	query := `INSERT INTO movies (name, date, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	result, err := db.Exec(query, m.Name, m.Date, m.UserID, m.CreatedAt, m.UpdatedAt)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err == nil {
		m.ID = int(id)
	}
	return err
}

func GetAllMovies(db *sql.DB) ([]Movie, error) {
	rows, err := db.Query(`SELECT id, name, date, user_id, created_at, updated_at FROM movies`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var m Movie
		if err := rows.Scan(&m.ID, &m.Name, &m.Date, &m.UserID, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, nil
}

func GetAllMoviesWithUserID(db *sql.DB, userID string) ([]Movie, error) {
	rows, err := db.Query(`SELECT id, name, date, user_id, created_at, updated_at FROM movies WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var m Movie
		if err := rows.Scan(&m.ID, &m.Name, &m.Date, &m.UserID, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, nil
}

func UpdateMovie(db *sql.DB, m Movie) error {
	query := `UPDATE movies SET name = ?, date = ?, user_id = ?, updated_at = ? WHERE id = ?`
	_, err := db.Exec(query, m.Name, m.Date, m.UserID, m.UpdatedAt, m.ID)
	return err
}

func DeleteMovie(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM movies WHERE id = ?`, id)
	return err
}
