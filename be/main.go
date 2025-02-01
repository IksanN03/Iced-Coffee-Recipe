package main

import (
	"be-test/database"
	"be-test/route"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	database.Init()

	// Set up Gin router
	router := gin.Default()

	// Set up routes
	route.SetupRoutes(router)

	// Run the server
	router.Run(":" + os.Getenv("PORT"))
}
