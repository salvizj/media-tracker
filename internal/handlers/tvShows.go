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

func TVShowsHandler(db *sql.DB, tmpl *template.Template) gin.HandlerFunc {
	return func(c *gin.Context) {
		tvShows, err := models.GetAllTVShows(db)
		if err != nil {
			c.String(http.StatusInternalServerError, "Neizdevās ielādēt TV šovus")
			return
		}

		data := types.LayoutTmplData{
			Title:           "TV shows",
			ContentTemplate: "content_tv_shows",
			TvShows:         tvShows,
		}
		c.HTML(http.StatusOK, "layout", data)
	}
}

func CreateTVShow(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tvShow models.TVShow
		if err := c.BindJSON(&tvShow); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
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
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		if err := models.DeleteTVShow(db, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete record"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
