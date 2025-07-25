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
	r.POST("/login", handlers.LoginHandler(db, tmpl))
	r.GET("/register", handlers.RegisterHandler(db, tmpl))
	r.POST("/register", handlers.RegisterHandler(db, tmpl))
	r.POST("/logout", handlers.LogoutHandler(db, tmpl))

	r.NoRoute(handlers.NotFoundHandler(db, tmpl))
	r.NoMethod(handlers.NotFoundHandler(db, tmpl))

	securedPages := r.Group("/")
	securedPages.Use(middleware.AuthRequired(db))
	{
		securedPages.GET("/manhwa-and-manga", handlers.ManhwaAndMangaHandler(db, tmpl))
		securedPages.GET("/movies", handlers.MoviesHandler(db, tmpl))
		securedPages.GET("/tv-shows", handlers.TVShowsHandler(db, tmpl))

		securedPages.GET("/download/tv-shows", handlers.DownloadTVShowsHandler(db))
		securedPages.GET("/download/manhwa-and-manga", handlers.DownloadManhwaAndMangaHandler(db))
		securedPages.GET("/download/movies", handlers.DownloadMoviesHandler(db))
		// Bulk add endpoints
		securedPages.POST("/bulk-add/tv-shows", handlers.BulkAddTVShowsHandler(db))
		securedPages.POST("/bulk-add/manhwa-and-manga", handlers.BulkAddManhwaAndMangaHandler(db))
		securedPages.POST("/bulk-add/movies", handlers.BulkAddMoviesHandler(db))
	}

	api := r.Group("/api")

	securedApi := api.Group("")
	securedApi.Use(middleware.AuthRequired(db))
	{
		securedApi.POST("/movies", handlers.CreateMovie(db))
		securedApi.PUT("/movies/:id", handlers.UpdateMovie(db))
		securedApi.DELETE("/movies/:id", handlers.DeleteMovie(db))

		securedApi.POST("/tv-shows", handlers.CreateTVShow(db))
		securedApi.PUT("/tv-shows/:id", handlers.UpdateTVShow(db))
		securedApi.DELETE("/tv-shows/:id", handlers.DeleteTVShow(db))

		securedApi.POST("/manhwa-and-manga", handlers.CreateManhwaAndManga(db))
		securedApi.PUT("/manhwa-and-manga/:id", handlers.UpdateManhwaAndManga(db))
		securedApi.DELETE("/manhwa-and-manga/:id", handlers.DeleteManhwaAndManga(db))

	}

	return r
}
