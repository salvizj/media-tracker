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

func ManhwaAndMangaHandler(db *sql.DB, tmpl *template.Template) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		sessionID, exists := c.Get("session_id")
		validSession, err := models.IsSessionValid(db, sessionID.(string))
		if !exists || err != nil || !validSession {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		entries, err := models.GetAllManhwaAndManga(db, userID.(string))
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to load manhwa and manga")
			return
		}

		data := types.LayoutTmplData{
			Title:           "Manhwa un Manga",
			ContentTemplate: "content_manhwa_and_manga",
			ManhwaAndManga:  entries,
		}
		c.HTML(http.StatusOK, "layout", data)
	}
}

func CreateManhwaAndManga(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var m models.ManhwaAndMangaInput
		if err := c.BindJSON(&m); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		manhwaAndManga := models.ManhwaAndManga{
			UserID:    m.UserID,
			Name:      m.Name,
			Status:    models.Status(m.Status),
			Chapter:   m.Chapter,
			Date:      time.Now().Format("2006-01-02"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := models.InsertManhwaAndManga(db, &manhwaAndManga); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entry"})
			return
		}
		c.JSON(http.StatusCreated, m)
	}
}

func UpdateManhwaAndManga(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var m models.ManhwaAndManga
		if err := c.BindJSON(&m); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		m.UpdatedAt = time.Now()
		if err := models.UpdateManhwaAndManga(db, m); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update entry"})
			return
		}
		c.JSON(http.StatusOK, m)
	}
}

func DeleteManhwaAndManga(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
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

		if err := models.DeleteManhwaAndManga(db, id, userID.(string)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete entry"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func DownloadManhwaAndMangaHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		entries, err := models.GetAllManhwaAndManga(db, userID.(string))
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to load manhwa and manga")
			return
		}
		c.Header("Content-Disposition", "attachment; filename=manhwa_and_manga.txt")
		c.Header("Content-Type", "text/plain")
		for _, entry := range entries {
			line := entry.Name + " | Status: " + string(entry.Status) + " | Chapter: " + strconv.Itoa(entry.Chapter) + " | Date: " + entry.Date + "\n"
			c.Writer.WriteString(line)
		}
	}
}

func BulkAddManhwaAndMangaHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		var items []struct {
			Name    string `json:"name"`
			Status  string `json:"status"`
			Chapter int    `json:"chapter"`
			Date    string `json:"date"`
		}
		if err := c.ShouldBindJSON(&items); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
			return
		}
		added, errors := 0, 0
		for _, item := range items {
			if item.Name == "" || item.Status == "" || item.Chapter < 1 {
				errors++
				continue
			}
			entry := models.ManhwaAndManga{
				Name:      item.Name,
				Status:    models.Status(item.Status),
				Chapter:   item.Chapter,
				Date:      item.Date,
				UserID:    userID.(string),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := models.InsertManhwaAndManga(db, &entry); err == nil {
				added++
			} else {
				errors++
			}
		}
		msg := "Added " + strconv.Itoa(added) + " manhwa/manga."
		if errors > 0 {
			msg += " " + strconv.Itoa(errors) + " lines had errors."
		}
		c.JSON(http.StatusOK, gin.H{"message": msg})
	}
}
