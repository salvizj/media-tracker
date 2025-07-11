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
		}

		if c.Request.Method == http.MethodPost {
			var user models.User
			if err := c.BindJSON(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
				return
			}
			userID, err := models.LoginUser(db, &user)
			if err != nil {
				c.HTML(http.StatusInternalServerError, "layout", types.LayoutTmplData{
					Title:           "Media Tracker",
					ContentTemplate: "content_login",
					Error:           "Invalid email or password",
				})
				return
			}
			user.ID = userID
			session := models.Session{
				ID:        utils.GenerateID(),
				UserID:    user.ID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			err = models.InsertSession(db, &session)
			if err != nil {
				c.HTML(http.StatusInternalServerError, "layout", types.LayoutTmplData{
					Title:           "Media Tracker",
					ContentTemplate: "content_login",
					Error:           "Failed to create session",
				})
				return
			}

			c.SetCookie("session_id", session.ID, 3600, "/", "", false, true)

			c.HTML(http.StatusOK, "layout", types.LayoutTmplData{
				Title:           "Media Tracker",
				ContentTemplate: "content_index",
				IsLoggedIn:      true,
			})
			return
		}
	}
}
