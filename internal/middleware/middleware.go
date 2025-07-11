package middleware

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthRequired(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("session_token")
		if err != nil || cookie == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var userID int
		err = db.QueryRow("SELECT user_id FROM sessions WHERE token = ? LIMIT 1", cookie).Scan(&userID)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
