package handlers

import (
	"database/sql"
	"html/template"
	"media_tracker/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LogoutHandler(db *sql.DB, tmpl *template.Template) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("session_id")
		if err != nil || cookie == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No active session"})
			return
		}
		sessionID := cookie
		models.DeleteSessionBySessionID(db, sessionID)

		c.SetCookie("session_id", "", -1, "/", "", false, true)
		c.SetCookie("user_id", "", -1, "/", "", false, true)

		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	}
}
