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
		entries, err := models.GetAllManhwaAndManga(db)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to load manhwa and manga")
			return
		}

		isLoggedIn, exists := c.Get("isLoggedIn")
		if !exists {
			isLoggedIn = false
		}

		data := types.LayoutTmplData{
			Title:           "Manhwa un Manga",
			ContentTemplate: "content_manhwa_and_manga",
			ManhwaAndManga:  entries,
			IsLoggedIn:      isLoggedIn.(bool),
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

		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		m.ID = id
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

		if err := models.DeleteManhwaAndManga(db, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete entry"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
