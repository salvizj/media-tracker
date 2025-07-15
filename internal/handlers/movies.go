package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"media_tracker/internal/models"
	"media_tracker/internal/types"

	"github.com/gin-gonic/gin"
)

func MoviesHandler(db *sql.DB, tmpl *template.Template) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.Redirect(http.StatusFound, "/login")
			return
		}
		movies, err := models.GetAllMovies(db, userID.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load movies"})
			return
		}

		data := types.LayoutTmplData{
			Title:           "Films",
			ContentTemplate: "content_movies",
			Movies:          movies,
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
			Date:      time.Now().Format("2006-01-02"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      m.Name,
			UserID:    m.UserID,
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
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		sessionID, exists := c.Get("session_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		validSession, err := models.IsSessionValid(db, sessionID.(string))
		if !exists || err != nil || !validSession {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		if err := models.DeleteMovie(db, id, userID.(string)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func DownloadMoviesHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		movies, err := models.GetAllMovies(db, userID.(string))
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to load movies")
			return
		}
		c.Header("Content-Disposition", "attachment; filename=movies.txt")
		c.Header("Content-Type", "text/plain")
		for _, movie := range movies {
			line := movie.Name + " | Date: " + movie.Date + "\n"
			c.Writer.WriteString(line)
		}
	}
}

func BulkAddMoviesHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		var items []struct {
			Name string `json:"name"`
			Date string `json:"date"`
		}
		if err := c.ShouldBindJSON(&items); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
			return
		}
		added, errors := 0, 0
		for _, item := range items {
			if item.Name == "" {
				errors++
				continue
			}
			movie := models.Movie{
				Name:      item.Name,
				Date:      item.Date,
				UserID:    userID.(string),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := models.InsertMovie(db, &movie); err == nil {
				added++
			} else {
				errors++
			}
		}
		msg := "Added " + strconv.Itoa(added) + " movies."
		if errors > 0 {
			msg += " " + strconv.Itoa(errors) + " lines had errors."
		}
		c.JSON(http.StatusOK, gin.H{"message": msg})
	}
}
