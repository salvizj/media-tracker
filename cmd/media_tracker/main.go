package main

import (
	"fmt"
	"html/template"
	"log"
	"media_tracker/internal/router"
	"media_tracker/internal/storage"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

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
	fmt.Printf("Started on http://localhost:%s\n", port)
	err = http.ListenAndServe(":"+port, router.NewRouter(db.DB, tmpl))
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}

}
