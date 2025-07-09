package models

import (
	"database/sql"
)

type ManhwaAndManga struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Date    string `json:"date"`
	Chapter int    `json:"chapter"`
}

type ManhwaAndMangaInput struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Chapter int    `json:"chapter"`
}

func CreateManhwaAndMangaTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS manhwa_and_manga (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		status TEXT,
		date TEXT,
		chapter INTEGER
	);`
	_, err := db.Exec(query)
	return err
}

func InsertManhwaAndManga(db *sql.DB, m *ManhwaAndManga) error {
	query := `INSERT INTO manhwa_and_manga (name, status, date, chapter) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, m.Name, m.Status, m.Date, m.Chapter)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err == nil {
		m.ID = int(id)
	}
	return err
}

func GetAllManhwaAndManga(db *sql.DB) ([]ManhwaAndManga, error) {
	rows, err := db.Query(`SELECT id, name, status, date, chapter FROM manhwa_and_manga`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []ManhwaAndManga
	for rows.Next() {
		var m ManhwaAndManga
		if err := rows.Scan(&m.ID, &m.Name, &m.Status, &m.Date, &m.Chapter); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

func UpdateManhwaAndManga(db *sql.DB, m ManhwaAndManga) error {
	query := `UPDATE manhwa_and_manga SET name = ?, status = ?, date = ?, chapter = ? WHERE id = ?`
	_, err := db.Exec(query, m.Name, m.Status, m.Date, m.Chapter, m.ID)
	return err
}

func DeleteManhwaAndManga(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM manhwa_and_manga WHERE id = ?`, id)
	return err
}
