package handlers

import (
	"database/sql"
	"html/template"
	"media_tracker/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotFoundHandler(db *sql.DB, tmpl *template.Template) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := types.LayoutTmplData{
			Title:           "Media Tracker",
			ContentTemplate: "404.html",
		}
		c.HTML(http.StatusNotFound, "layout", data)
	}
}
