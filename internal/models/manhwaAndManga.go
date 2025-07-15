package models

import (
	"database/sql"
	"time"
)

type Status string

const (
	StatusReading   Status = "Reading"
	StatusWatching  Status = "Watching"
	StatusCompleted Status = "Completed"
	StatusDropped   Status = "Dropped"
)

type ManhwaAndManga struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Status    Status    `json:"status"`
	Date      string    `json:"date"`
	Chapter   int       `json:"chapter"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    string    `json:"user_id"`
}

type ManhwaAndMangaInput struct {
	Name    string `json:"name"`
	Status  Status `json:"status"`
	Chapter int    `json:"chapter"`
	UserID  string `json:"user_id"`
}

func CreateManhwaAndMangaTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS manhwa_and_manga (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		status TEXT,
		date TEXT,
		chapter INTEGER,
		user_id TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(query)
	return err
}

func InsertManhwaAndManga(db *sql.DB, m *ManhwaAndManga) error {
	query := `INSERT INTO manhwa_and_manga (name, status, date, chapter, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, m.Name, m.Status, m.Date, m.Chapter, m.UserID, m.CreatedAt, m.UpdatedAt)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err == nil {
		m.ID = int(id)
	}
	return err
}

func GetAllManhwaAndManga(db *sql.DB, userID string) ([]ManhwaAndManga, error) {
	rows, err := db.Query(`SELECT id, name, status, date, chapter, user_id, created_at, updated_at FROM manhwa_and_manga WHERE user_id = ? ORDER BY date DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []ManhwaAndManga
	for rows.Next() {
		var m ManhwaAndManga
		if err := rows.Scan(&m.ID, &m.Name, &m.Status, &m.Date, &m.Chapter, &m.UserID, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

func GetAllManhwasAndMangas(db *sql.DB, userID string) ([]ManhwaAndManga, error) {
	rows, err := db.Query(`SELECT id, name, status, date, chapter, user_id, created_at, updated_at FROM manhwa_and_manga WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []ManhwaAndManga
	for rows.Next() {
		var m ManhwaAndManga
		if err := rows.Scan(&m.ID, &m.Name, &m.Status, &m.Date, &m.Chapter, &m.UserID, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

func UpdateManhwaAndManga(db *sql.DB, m ManhwaAndManga) error {
	query := `UPDATE manhwa_and_manga SET name = ?, status = ?, date = ?, chapter = ?, user_id = ?, updated_at = ? WHERE id = ? AND user_id = ?`
	_, err := db.Exec(query, m.Name, m.Status, m.Date, m.Chapter, m.UserID, m.UpdatedAt, m.ID, m.UserID)
	return err
}

func DeleteManhwaAndManga(db *sql.DB, id int, userID string) error {
	_, err := db.Exec(`DELETE FROM manhwa_and_manga WHERE id = ? AND user_id = ?`, id, userID)
	return err
}
