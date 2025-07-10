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

		data := types.LayoutTmplData{
			Title:           "Media Tracker",
			ContentTemplate: "content_index",
		}
		c.HTML(http.StatusOK, "layout", data)
	}
}
