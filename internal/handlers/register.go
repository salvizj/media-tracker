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

func RegisterHandler(db *sql.DB, tmpl *template.Template) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {

			data := types.LayoutTmplData{
				Title:           "Media Tracker",
				ContentTemplate: "content_register",
				IsLoggedIn:      false,
			}
			c.HTML(http.StatusOK, "layout", data)
		}

		if c.Request.Method == http.MethodPost {
			var user models.User
			if err := c.BindJSON(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
				return
			}
			if user.Password != user.ConfirmPassword {
				c.HTML(http.StatusBadRequest, "layout", types.LayoutTmplData{
					Title:           "Media Tracker",
					ContentTemplate: "content_register",
					Error:           "Passwords do not match",
					IsLoggedIn:      false,
				})
				return
			}
			var userInsert = models.User{
				ID:        utils.GenerateUUID(),
				Email:     user.Email,
				Password:  utils.HashPassword(user.Password),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			err := models.InsertUser(db, &userInsert)
			if err != nil {
				c.HTML(http.StatusInternalServerError, "layout", types.LayoutTmplData{
					Title:           "Media Tracker",
					ContentTemplate: "content_register",
					Error:           "Failed to create user",
					IsLoggedIn:      false,
				})
				return
			}
			c.HTML(http.StatusOK, "layout", types.LayoutTmplData{
				Title:           "Media Tracker",
				ContentTemplate: "content_login",
				IsLoggedIn:      false,
			})
			return
		}
	}
}
