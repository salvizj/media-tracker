package router

import (
	"database/sql"
	"html/template"
	"media_tracker/internal/handlers"
	"media_tracker/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(db *sql.DB, tmpl *template.Template) http.Handler {
	r := gin.Default()
	r.SetHTMLTemplate(tmpl)

	r.Static("/public", "./public")

	r.GET("/", handlers.IndexHandler(db, tmpl))
	r.GET("/login", handlers.LoginHandler(db, tmpl))
	r.GET("/register", handlers.RegisterHandler(db, tmpl))

	secured := r.Group("/")
	secured.Use(middleware.AuthRequired(db))
	{
		secured.GET("/manhwa-and-manga", handlers.ManhwaAndMangaHandler(db, tmpl))
		secured.GET("/movies", handlers.MoviesHandler(db, tmpl))
		secured.GET("/tv-shows", handlers.TVShowsHandler(db, tmpl))
	}

	api := r.Group("/api")
	api.Use(middleware.AuthRequired(db))
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
