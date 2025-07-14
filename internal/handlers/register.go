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

			if user.Password != user.ConfirmPassword {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
				return
			}

			var userInsert = models.User{
				ID:        utils.GenerateUUID(),
				Email:     user.Email,
				Password:  utils.HashPassword(user.Password),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			err, exists := models.InsertUser(db, &userInsert)
			if err != nil {
				if exists {
					c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Account created successfully! You can now log in."})
			return
		}
	}
}
