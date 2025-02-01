package handler

import (
	"be-test/database"
	"be-test/helpers"
	"be-test/models"
	"be-test/utils"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// SubmitEmail generates a magic link for email authentication
func SubmitEmail(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		helpers.NewAPIResponse(c, nil, err, "binding", 0, "Invalid email format")
		return
	}

	token, err := helpers.GenerateJWT(user.Email, time.Now().Add(time.Minute*5).Unix())
	if err != nil {
		helpers.NewAPIResponse(c, nil, err, "token", 0, "Failed to generate token")
		return
	}

	magicLink := os.Getenv("FRONTEND_URL") + "?token=" + token
	err = utils.SendEmail(user.Email, magicLink)
	if err != nil {
		helpers.NewAPIResponse(c, nil, err, "send_email", 0, "")

		return
	}

	helpers.NewAPIResponse(c, nil, nil, "", 0, "Magic link sent")
}

// MagicLink authenticates a user via the magic link
func MagicLink(c *gin.Context) {
	token := c.DefaultQuery("token", "")
	if token == "" {
		helpers.NewAPIResponse(c, nil, nil, "token", http.StatusBadRequest, "Token is required")
		return
	}

	// Validate JWT token
	claims, err := helpers.ValidateJWT(token)
	if err != nil {
		helpers.NewAPIResponse(c, nil, nil, "token", http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	email := claims["email"].(string)

	// Check if token already exists in database
	var existingToken models.User
	result := database.DB.Where("access_token = ?", token).First(&existingToken)
	if result.Error == nil {
		helpers.NewAPIResponse(c, nil, nil, "token", http.StatusUnauthorized, "Token already used")
		return
	}

	// Generate new access token
	newToken, err := helpers.GenerateJWT(email, time.Now().Add(time.Hour*24).Unix())
	if err != nil {
		helpers.NewAPIResponse(c, nil, nil, "newToken", 0, "Failed to generate token")
		return
	}

	// Check if user exists and update token, or create new user
	var user models.User
	result = database.DB.Where("email = ?", email).First(&user)
	if result.Error == nil {
		// User exists, update token
		database.DB.Model(&user).Update("access_token", token)
	} else {
		// Create new user
		user = models.User{
			Email:       email,
			AccessToken: token,
		}
		database.DB.Create(&user)
	}

	helpers.NewAPIResponse(c, gin.H{
		"access_token": newToken,
	}, nil, "authenticated", 0, "Authentication successful")

}
