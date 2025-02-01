package models

import "gorm.io/gorm"

// User represents a user in the system
type User struct {
	gorm.Model
	Email       string `json:"email"  binding:"required,email"`
	AccessToken string `json:"access_token"`
}
