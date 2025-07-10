package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"media_tracker/internal/router"
	"media_tracker/internal/storage"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mode := flag.String("mode", "debug", "set gin mode: debug or release")
	flag.Parse()
	gin.SetMode(*mode)

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
	err = http.ListenAndServe("0.0.0.0:"+port, router.NewRouter(db.DB, tmpl))
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}

}
