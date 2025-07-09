package models

import (
	"database/sql"
)

type TVShow struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Date    string `json:"date"`
	Season  int    `json:"season"`
	Episode int    `json:"episode"`
}

type TVShowInput struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Season  int    `json:"season"`
	Episode int    `json:"episode"`
}

func CreateTVShowsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS tv_shows (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		status TEXT,
		date TEXT,
		season INTEGER,
		episode INTEGER
	);`
	_, err := db.Exec(query)
	return err
}

func InsertTVShow(db *sql.DB, show *TVShow) error {
	query := `INSERT INTO tv_shows (name, status, date, season, episode) VALUES (?, ?, ?, ?, ?)`
	result, err := db.Exec(query, show.Name, show.Status, show.Date, show.Season, show.Episode)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err == nil {
		show.ID = int(id)
	}
	return err
}

func GetAllTVShows(db *sql.DB) ([]TVShow, error) {
	rows, err := db.Query(`SELECT id, name, status, date, season, episode FROM tv_shows`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shows []TVShow
	for rows.Next() {
		var s TVShow
		if err := rows.Scan(&s.ID, &s.Name, &s.Status, &s.Date, &s.Season, &s.Episode); err != nil {
			return nil, err
		}
		shows = append(shows, s)
	}
	return shows, nil
}

func UpdateTVShow(db *sql.DB, s TVShow) error {
	query := `UPDATE tv_shows SET name = ?, status = ?, date = ?, season = ?, episode = ? WHERE id = ?`
	_, err := db.Exec(query, s.Name, s.Status, s.Date, s.Season, s.Episode, s.ID)
	return err
}

func DeleteTVShow(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM tv_shows WHERE id = ?`, id)
	return err
}
