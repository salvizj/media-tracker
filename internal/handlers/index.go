package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"media_tracker/internal/types"

	"github.com/gin-gonic/gin"
)

func IndexHandler(db *sql.DB, tmpl *template.Template) gin.HandlerFunc {
	return func(c *gin.Context) {
		isLoggedIn, exists := c.Get("isLoggedIn")
		if !exists {
			isLoggedIn = false
		}

		data := types.LayoutTmplData{
			Title:           "Media Tracker",
			ContentTemplate: "content_index",
			IsLoggedIn:      isLoggedIn.(bool),
		}
		c.HTML(http.StatusOK, "layout", data)
	}
}
