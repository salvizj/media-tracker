package models

import (
	"database/sql"
	"time"
)

type Movie struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Date time.Time
}

func CreateMoviesTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS movies (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		date TEXT
	);`
	_, err := db.Exec(query)
	return err
}

func InsertMovie(db *sql.DB, m *Movie) error {
	query := `INSERT INTO movies (name, date) VALUES (?, ?)`
	result, err := db.Exec(query, m.Name, m.Date)
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
	rows, err := db.Query(`SELECT id, name, date FROM movies`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var m Movie
		if err := rows.Scan(&m.ID, &m.Name, &m.Date); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, nil
}

func UpdateMovie(db *sql.DB, m Movie) error {
	query := `UPDATE movies SET name = ?, date = ? WHERE id = ?`
	_, err := db.Exec(query, m.Name, m.Date, m.ID)
	return err
}

func DeleteMovie(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM movies WHERE id = ?`, id)
	return err
}
