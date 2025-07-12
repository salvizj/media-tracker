package handlers

import (
	"database/sql"
	"html/template"
	"media_tracker/internal/models"
	"media_tracker/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LogoutHandler(db *sql.DB, tmpl *template.Template) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("session_token")
		if err != nil || cookie == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		sessionID := cookie
		models.DeleteSessionBySessionID(db, sessionID)
		c.HTML(http.StatusOK, "layout", types.LayoutTmplData{
			Title:           "Media Tracker",
			ContentTemplate: "content_logout",
		})

	}
}
