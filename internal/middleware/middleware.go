package middleware

import (
	"database/sql"
	"media_tracker/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthRequired(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("session_id")
		if err != nil || cookie == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		session, err := models.GetSessionBySessionID(db, cookie)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("userID", session.UserID)
		c.Next()
	}
}
