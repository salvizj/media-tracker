package storage

import (
	"database/sql"
	"log"
	"media_tracker/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	DB *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		db.Close()
		return nil, err
	}

	return &Storage{DB: db}, nil
}

func createTables(db *sql.DB) error {
	if err := models.CreateMoviesTable(db); err != nil {
		return err
	}
	if err := models.CreateTVShowsTable(db); err != nil {
		return err
	}
	if err := models.CreateManhwaAndMangaTable(db); err != nil {
		return err
	}
	return nil
}

func (s *Storage) Close() {
	if err := s.DB.Close(); err != nil {
		log.Println("Error closing DB:", err)
	}
}
