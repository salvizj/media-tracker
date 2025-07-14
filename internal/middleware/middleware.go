package middleware

import (
	"database/sql"
	"media_tracker/internal/models"
	"media_tracker/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthRequired(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieValue, err := c.Cookie("session_id")
		if err != nil || cookieValue == "" {
			data := types.LayoutTmplData{
				Title:           "Media Tracker",
				ContentTemplate: "content_restricted",
			}
			c.HTML(http.StatusUnauthorized, "layout", data)
			return
		}
		userID, err := models.GetUserIDBySessionID(db, cookieValue)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		valid, err := models.IsSessionValid(db, cookieValue)
		if err != nil || !valid {
			data := types.LayoutTmplData{
				Title:           "Media Tracker",
				ContentTemplate: "content_restricted",
			}
			c.HTML(http.StatusUnauthorized, "layout", data)
			return
		}
		c.Set("session_id", cookieValue)
		c.Set("user_id", userID)
		c.Next()
	}
}
