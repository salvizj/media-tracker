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
			isLoggedIn, exists := c.Get("isLoggedIn")
			if !exists {
				isLoggedIn = false
			}
			data := types.LayoutTmplData{
				Title:           "Media Tracker",
				ContentTemplate: "content_register",
				IsLoggedIn:      isLoggedIn.(bool),
			}
			c.HTML(http.StatusOK, "layout", data)
		}

		if c.Request.Method == http.MethodPost {
			var user models.User
			if err := c.BindJSON(&user); err != nil {
				if c.Request.URL.Path == "/api/register" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
				} else {
					data := types.LayoutTmplData{
						Title:           "Media Tracker",
						ContentTemplate: "content_register",
						IsLoggedIn:      false,
					}
					data.Error = "Invalid input"
					c.HTML(http.StatusBadRequest, "layout", data)
				}
				return
			}

			if user.Password != user.ConfirmPassword {
				if c.Request.URL.Path == "/api/register" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
				} else {

					data := types.LayoutTmplData{
						Title:           "Media Tracker",
						ContentTemplate: "content_register",
						IsLoggedIn:      false,
					}
					data.Error = "Passwords do not match"
					c.HTML(http.StatusBadRequest, "layout", data)
				}
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
				if c.Request.URL.Path == "/api/register" {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				} else {
					data := types.LayoutTmplData{
						Title:           "Media Tracker",
						ContentTemplate: "content_register",
						IsLoggedIn:      false,
					}
					data.Error = "Failed to create user"
					c.HTML(http.StatusInternalServerError, "layout", data)
				}
				return
			}

			if c.Request.URL.Path == "/api/register" {
				c.JSON(http.StatusCreated, gin.H{
					"message": "Account created successfully! You can now log in.",
					"userID":  userInsert.ID,
				})
			} else {
				data := types.LayoutTmplData{
					Title:           "Media Tracker",
					ContentTemplate: "content_login",
					IsLoggedIn:      false,
				}
				c.HTML(http.StatusOK, "layout", data)
			}
		}
	}
}
