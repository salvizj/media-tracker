package handlers

import (
	"database/sql"
	"html/template"
	"media_tracker/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(db *sql.DB, tmpl *template.Template) gin.HandlerFunc {
	return func(c *gin.Context) {

		data := types.LayoutTmplData{
			Title:           "Media Tracker",
			ContentTemplate: "content_register",
		}
		c.HTML(http.StatusOK, "layout", data)
	}

}
