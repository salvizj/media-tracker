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

	"github.com/gin-gonic/gin"
)

func ManhwaAndMangaHandler(db *sql.DB, tmpl *template.Template) gin.HandlerFunc {
	return func(c *gin.Context) {
<<<<<<< HEAD
		userID, err := c.Cookie("user_id")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
			return
		}
		entries, err := models.GetAllManhwasAndMangasWithUserID(db, userID)
=======
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
>>>>>>> feature/auth
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
			Status:    m.Status,
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

func DownloadManhwasAndMangas(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		manhwaAndManga, err := models.GetAllManhwasAndMangas(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load manhwa and manga"})
			return
		}

		var builder strings.Builder
		for _, manhwaAndManga := range manhwaAndManga {
			builder.WriteString(manhwaAndManga.Name + "\n")
			builder.WriteString(manhwaAndManga.Date + "\n")
			builder.WriteString(manhwaAndManga.CreatedAt.Format("2006-01-02") + "\n")
			builder.WriteString(manhwaAndManga.UpdatedAt.Format("2006-01-02") + "\n")
			builder.WriteString(manhwaAndManga.UserID + "\n")
		}
		content := builder.String()

		filename := "manhwa_and_manga.txt"
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(content))
	}
}
