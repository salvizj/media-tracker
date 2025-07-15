package models

import (
	"database/sql"
	"time"
)

type TVShow struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Status    Status    `json:"status"`
	Date      string    `json:"date"`
	Season    int       `json:"season"`
	Episode   int       `json:"episode"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    string    `json:"user_id"`
}

type TVShowInput struct {
	Name    string `json:"name"`
	Status  Status `json:"status"`
	Season  int    `json:"season"`
	Episode int    `json:"episode"`
	UserID  string `json:"user_id"`
}

func CreateTVShowsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS tv_shows (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		status TEXT,
		date TEXT,
		season INTEGER,
		episode INTEGER,
		user_id TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(query)
	return err
}

func InsertTVShow(db *sql.DB, show *TVShow) error {
	query := `INSERT INTO tv_shows (name, status, date, season, episode, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, show.Name, show.Status, show.Date, show.Season, show.Episode, show.UserID, show.CreatedAt, show.UpdatedAt, show.UserID)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err == nil {
		show.ID = int(id)
	}
	return err
}

func GetAllTVShows(db *sql.DB, userID string) ([]TVShow, error) {
	rows, err := db.Query(`SELECT id, name, status, date, season, episode, user_id, created_at, updated_at FROM tv_shows WHERE user_id = ? ORDER BY date DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shows []TVShow
	for rows.Next() {
		var s TVShow
		if err := rows.Scan(&s.ID, &s.Name, &s.Status, &s.Date, &s.Season, &s.Episode, &s.UserID, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		shows = append(shows, s)
	}
	return shows, nil
}
func GetAllTVShowsWithUserID(db *sql.DB, userID string) ([]TVShow, error) {
	rows, err := db.Query(`SELECT id, name, status, date, season, episode, user_id, created_at, updated_at FROM tv_shows WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shows []TVShow
	for rows.Next() {
		var s TVShow
		if err := rows.Scan(&s.ID, &s.Name, &s.Status, &s.Date, &s.Season, &s.Episode, &s.UserID, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		shows = append(shows, s)
	}
	return shows, nil
}

func UpdateTVShow(db *sql.DB, s TVShow) error {
	query := `UPDATE tv_shows SET name = ?, status = ?, date = ?, season = ?, episode = ?, user_id = ?, updated_at = ? WHERE id = ? AND user_id = ?`
	_, err := db.Exec(query, s.Name, s.Status, s.Date, s.Season, s.Episode, s.UserID, s.UpdatedAt, s.ID, s.UserID)
	return err
}

func DeleteTVShow(db *sql.DB, id int, userID string) error {
	_, err := db.Exec(`DELETE FROM tv_shows WHERE id = ? AND user_id = ?`, id, userID)
	return err
}
