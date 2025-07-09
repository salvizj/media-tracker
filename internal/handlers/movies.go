package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"media_tracker/internal/models"
	"media_tracker/internal/types"

	"log"

	"github.com/gin-gonic/gin"
)

func MoviesHandler(db *sql.DB, tmpl *template.Template) gin.HandlerFunc {
	return func(c *gin.Context) {
		movies, err := models.GetAllMovies(db)
		if err != nil {
			log.Println("Error loading movies:", err)
			c.String(http.StatusInternalServerError, "Failed to load movies")
			return
		}

		data := types.LayoutTmplData{
			Title:           "Filmas",
			ContentTemplate: "content_movies",
			Movies:          movies,
		}
		c.HTML(http.StatusOK, "layout", data)
	}
}

func CreateMovie(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var m models.Movie
		if err := c.BindJSON(&m); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		m.Date = time.Now().Format("2006-01-02")
		if err := models.InsertMovie(db, &m); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create movie"})
			return
		}
		c.JSON(http.StatusCreated, m)
	}
}

func UpdateMovie(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var m models.Movie
		if err := c.BindJSON(&m); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		m.ID = id

		if err := models.UpdateMovie(db, m); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
			return
		}
		c.JSON(http.StatusOK, m)
	}
}

func DeleteMovie(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		if err := models.DeleteMovie(db, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
