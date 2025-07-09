package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"

	"media_tracker/internal/models"
	"media_tracker/internal/types"

	"github.com/gin-gonic/gin"
)

func ManhwaAndMangaHandler(db *sql.DB, tmpl *template.Template) gin.HandlerFunc {
	return func(c *gin.Context) {
		entries, err := models.GetAllManhwaAndManga(db)
		if err != nil {
			c.String(http.StatusInternalServerError, "Neizdevās ielādēt manhu un manga")
			return
		}

		data := types.LayoutTmplData{
			Title:           "Manhwa un Manga",
			Message:         "Labākās manhu un manga sērijas!",
			ContentTemplate: "content_manhwa_and_manga",
			ManhwaAndManga:  entries,
		}
		c.HTML(http.StatusOK, "layout", data)
	}
}

func CreateManhwaAndManga(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var m models.ManhwaAndManga
		if err := c.BindJSON(&m); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nederīga ievade"})
			return
		}
		if err := models.InsertManhwaAndManga(db, &m); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Neizdevās izveidot ierakstu"})
			return
		}
		c.JSON(http.StatusCreated, m)
	}
}

func UpdateManhwaAndManga(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var m models.ManhwaAndManga
		if err := c.BindJSON(&m); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nederīga ievade"})
			return
		}

		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nederīgs ID"})
			return
		}
		m.ID = id

		if err := models.UpdateManhwaAndManga(db, m); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Neizdevās atjaunināt ierakstu"})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nederīgs ID"})
			return
		}

		if err := models.DeleteManhwaAndManga(db, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Neizdevās dzēst ierakstu"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
