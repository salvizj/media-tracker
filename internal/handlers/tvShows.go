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

func TVShowsHandler(db *sql.DB, tmpl *template.Template) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		tvShows, err := models.GetAllTVShows(db, userID.(string))
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to load tv shows")
			return
		}
		data := types.LayoutTmplData{
			Title:           "TV shows",
			ContentTemplate: "content_tv_shows",
			TVShows:         tvShows,
		}
		c.HTML(http.StatusOK, "layout", data)
	}
}

func CreateTVShow(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.TVShowInput
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		tvShow := models.TVShow{
			UserID:    input.UserID,
			Name:      input.Name,
			Status:    input.Status,
			Season:    input.Season,
			Episode:   input.Episode,
			Date:      time.Now().Format("2006-01-02"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := models.InsertTVShow(db, &tvShow); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create record"})
			return
		}
		c.JSON(http.StatusCreated, tvShow)
	}
}

func UpdateTVShow(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tvShow models.TVShow
		if err := c.BindJSON(&tvShow); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		tvShow.UpdatedAt = time.Now()
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		tvShow.ID = id

		if err := models.UpdateTVShow(db, tvShow); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
			return
		}
		c.JSON(http.StatusOK, tvShow)
	}
}

func DeleteTVShow(db *sql.DB) gin.HandlerFunc {
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

		if err := models.DeleteTVShow(db, id, userID.(string)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete record"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func DownloadTVShowsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		tvShows, err := models.GetAllTVShows(db, userID.(string))
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to load tv shows")
			return
		}
		c.Header("Content-Disposition", "attachment; filename=tv_shows.txt")
		c.Header("Content-Type", "text/plain")
		for _, show := range tvShows {
			line := show.Name + " | Status: " + string(show.Status) + " | Season: " + strconv.Itoa(show.Season) + " | Episode: " + strconv.Itoa(show.Episode) + " | Date: " + show.Date + "\n"
			c.Writer.WriteString(line)
		}
	}
}

func BulkAddTVShowsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		var items []struct {
			Name    string `json:"name"`
			Status  string `json:"status"`
			Season  int    `json:"season"`
			Episode int    `json:"episode"`
			Date    string `json:"date"`
		}
		if err := c.ShouldBindJSON(&items); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
			return
		}
		added, errors := 0, 0
		for _, item := range items {
			if item.Name == "" || item.Status == "" || item.Season < 1 || item.Episode < 1 {
				errors++
				continue
			}
			tvShow := models.TVShow{
				UserID:    userID.(string),
				Name:      item.Name,
				Status:    models.Status(item.Status),
				Season:    item.Season,
				Episode:   item.Episode,
				Date:      item.Date,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := models.InsertTVShow(db, &tvShow); err == nil {
				added++
			} else {
				errors++
			}
		}
		msg := "Added " + strconv.Itoa(added) + " TV shows."
		if errors > 0 {
			msg += " " + strconv.Itoa(errors) + " lines had errors."
		}
		c.JSON(http.StatusOK, gin.H{"message": msg})
	}
}
