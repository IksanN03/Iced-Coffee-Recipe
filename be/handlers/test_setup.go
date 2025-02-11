package handler

import (
	"be-test/database"
	"be-test/middleware"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	// Initialize database using existing setup
	database.InitTest()

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Auth routes
	r.POST("/auth/submit-email", SubmitEmail)
	r.GET("/auth/magic-link", MagicLink)

	// Protected routes
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		// Inventory routes
		authorized.POST("/inventory", AddInventory)
		authorized.GET("/inventory", GetInventory)
		authorized.PUT("/inventory/:id", UpdateInventory)
		authorized.DELETE("/inventory/:id", DeleteInventory)

		// Recipe routes
		authorized.POST("/recipe", AddRecipe)
		authorized.GET("/recipe", GetRecipe)
		authorized.PUT("/recipe/:id", UpdateRecipe)
	}

	return r
}
