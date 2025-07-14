package handlers

import (
	"database/sql"
	"html/template"
	"media_tracker/internal/models"
	"media_tracker/internal/types"
	"media_tracker/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func LoginHandler(db *sql.DB, tmpl *template.Template) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {
			data := types.LayoutTmplData{
				Title:           "Media Tracker",
				ContentTemplate: "content_login",
			}
			c.HTML(http.StatusOK, "layout", data)
			return
		}

		if c.Request.Method == http.MethodPost {
			var user models.User
			if err := c.BindJSON(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
				return
			}

			userID, err := models.LoginUser(db, &user)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
				return
			}

			session := models.Session{
				ID:        utils.GenerateUUID(),
				UserID:    userID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				ExpiresAt: time.Now().Add(1 * time.Hour),
			}

			err = models.InsertSession(db, &session)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
				return
			}

			_ = models.CleanUpExpiredSessions(db)

			sessionCookie := http.Cookie{
				Name:     "session_id",
				Value:    session.ID,
				MaxAge:   3600,
				Path:     "/",
				Domain:   "",
				HttpOnly: false,
				Secure:   false,
			}
			c.SetCookie(sessionCookie.Name, sessionCookie.Value, sessionCookie.MaxAge, sessionCookie.Path, sessionCookie.Domain, sessionCookie.Secure, sessionCookie.HttpOnly)

			userCookie := http.Cookie{
				Name:     "user_id",
				Value:    user.ID,
				MaxAge:   3600,
				Path:     "/",
				Domain:   "",
				HttpOnly: false,
				Secure:   false,
			}
			c.SetCookie(userCookie.Name, userCookie.Value, userCookie.MaxAge, userCookie.Path, userCookie.Domain, userCookie.Secure, userCookie.HttpOnly)

			c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
			return
		}
	}
}
