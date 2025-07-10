package router

import (
	"database/sql"
	"html/template"
	"media_tracker/internal/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(db *sql.DB, tmpl *template.Template) http.Handler {
	r := gin.Default()
	r.SetHTMLTemplate(tmpl)

	r.Static("/public", "./public")

	// UI routes
	r.GET("/manhwa-and-manga", handlers.ManhwaAndMangaHandler(db, tmpl))
	r.GET("/movies", handlers.MoviesHandler(db, tmpl))
	r.GET("/tv-shows", handlers.TVShowsHandler(db, tmpl))

	// CRUD API routes
	api := r.Group("/api")
	{
		api.POST("/movies", handlers.CreateMovie(db))
		api.PUT("/movies/:id", handlers.UpdateMovie(db))
		api.DELETE("/movies/:id", handlers.DeleteMovie(db))

		api.POST("/tv-shows", handlers.CreateTVShow(db))
		api.PUT("/tv-shows/:id", handlers.UpdateTVShow(db))
		api.DELETE("/tv-shows/:id", handlers.DeleteTVShow(db))

		api.POST("/manhwa-and-manga", handlers.CreateManhwaAndManga(db))
		api.PUT("/manhwa-and-manga/:id", handlers.UpdateManhwaAndManga(db))
		api.DELETE("/manhwa-and-manga/:id", handlers.DeleteManhwaAndManga(db))
	}

	return r
}
