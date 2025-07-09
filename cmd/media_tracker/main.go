package main

import (
	"fmt"
	"html/template"
	"log"
	"media_tracker/internal/router"
	"media_tracker/internal/storage"
	"net/http"
)

func main() {
	db, err := storage.New("./media_tracker.db")
	if err != nil {
		log.Fatalf("Cannot initialize storage: %v", err)
	}
	defer db.Close()
	tmpl, err := template.ParseFiles(
		"templates/layout.html",
		"templates/movies.html",
		"templates/tv_shows.html",
		"templates/manhwa_and_manga.html",
	)
	if err != nil {
		log.Fatalf("Cannot parse templates: %v", err)
	}
	fmt.Printf("Started on %s \n", "http://localhost:8080")
	http.ListenAndServe(":8080", router.NewRouter(db.DB, tmpl))
}
