package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"strings"
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

		isLoggedIn, exists := c.Get("isLoggedIn")
		if !exists {
			isLoggedIn = false
		}

		data := types.LayoutTmplData{
			Title:           "Films",
			ContentTemplate: "content_movies",
			Movies:          movies,
			IsLoggedIn:      isLoggedIn.(bool),
		}
		c.HTML(http.StatusOK, "layout", data)
	}
}

func CreateMovie(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var m models.MovieInput
		if err := c.BindJSON(&m); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		movie := models.Movie{
			Name:      m.Name,
			Date:      time.Now().Format("2006-01-02"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := models.InsertMovie(db, &movie); err != nil {
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
		m.UpdatedAt = time.Now()
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

func DownloadMovies(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := c.Cookie("user_id")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
			return
		}
		movies, err := models.GetAllMoviesWithUserID(db, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load movies"})
			return
		}

		var builder strings.Builder
		for _, movie := range movies {
			builder.WriteString(movie.Name + ";")
			builder.WriteString(movie.Date + ";")
			builder.WriteString(movie.CreatedAt.Format("2006-01-02") + ";")
			builder.WriteString(movie.UpdatedAt.Format("2006-01-02") + "\n")
		}
		content := builder.String()

		filename := "movies.txt"
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(content))
	}
}
