package route

import (
	handler "be-test/handlers"
	"be-test/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all the API routes
func SetupRoutes(router *gin.Engine) {

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Authentication Routes
	router.POST("/auth/submit-email", handler.SubmitEmail)
	router.GET("/auth/magic-link", handler.MagicLink)

	// Protected Routes (requires JWT authentication)
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())

	// Inventory Routes
	protected.GET("/inventory", handler.GetInventory)
	protected.POST("/inventory", handler.AddInventory)
	protected.PUT("/inventory/:id", handler.UpdateInventory)
	protected.DELETE("/inventory/:id", handler.DeleteInventory)

	// Recipe Routes
	protected.POST("/recipe", handler.AddRecipe)
	protected.GET("/recipe", handler.GetRecipe)
	protected.PUT("/recipe/:id", handler.UpdateRecipe)
}
