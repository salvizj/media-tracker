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
			isLoggedIn, exists := c.Get("isLoggedIn")
			if !exists {
				isLoggedIn = false
			}
			data := types.LayoutTmplData{
				Title:           "Media Tracker",
				ContentTemplate: "content_login",
				IsLoggedIn:      isLoggedIn.(bool),
			}
			c.HTML(http.StatusOK, "layout", data)
		}

		if c.Request.Method == http.MethodPost {
			var user models.User
			if err := c.BindJSON(&user); err != nil {
				if c.Request.URL.Path == "/api/login" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
				} else {
					data := types.LayoutTmplData{
						Title:           "Media Tracker",
						ContentTemplate: "content_login",
						Error:           "Invalid input",
						IsLoggedIn:      false,
					}
					data.Error = "Invalid input"
					c.HTML(http.StatusBadRequest, "layout", data)
				}
				return
			}

			userID, err := models.LoginUser(db, &user)
			if err != nil {
				if c.Request.URL.Path == "/api/login" {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
				} else {
					data := types.LayoutTmplData{
						Title:           "Media Tracker",
						ContentTemplate: "content_login",
						Error:           "Invalid email or password",
						IsLoggedIn:      false,
					}
					data.Error = "Invalid email or password"
					c.HTML(http.StatusInternalServerError, "layout", data)
				}
				return
			}

			user.ID = userID
			session := models.Session{
				ID:        utils.GenerateUUID(),
				UserID:    user.ID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				ExpiresAt: time.Now().Add(24 * time.Hour),
			}

			err = models.InsertSession(db, &session)
			if err != nil {
				if c.Request.URL.Path == "/api/login" {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
				} else {
					data := types.LayoutTmplData{
						Title:           "Media Tracker",
						ContentTemplate: "content_login",
						Error:           "Failed to create session",
						IsLoggedIn:      false,
					}
					data.Error = "Failed to create session"
					c.HTML(http.StatusInternalServerError, "layout", data)
				}
				return
			}

			models.CleanUpExpiredSessions(db)
			cookie := http.Cookie{
				Name:     "session_id",
				Value:    session.ID,
				Expires:  session.ExpiresAt,
				HttpOnly: true,
				Secure:   true,
			}
			c.SetCookie(cookie.Name, cookie.Value, 3600, "/", "", false, true)

			userCookie := http.Cookie{
				Name:     "user_id",
				Value:    user.ID,
				Expires:  time.Now().Add(24 * time.Hour),
				HttpOnly: false,
				Secure:   true,
			}
			c.SetCookie(userCookie.Name, userCookie.Value, 3600, "/", "", false, true)

			if c.Request.URL.Path == "/api/login" {
				c.JSON(http.StatusOK, gin.H{
					"message": "Login successful",
					"userID":  user.ID,
				})
			} else {
				data := types.LayoutTmplData{
					Title:           "Media Tracker",
					ContentTemplate: "content_index",
					IsLoggedIn:      true,
				}
				c.HTML(http.StatusOK, "layout", data)
			}
		}
	}
}
