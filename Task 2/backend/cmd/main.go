package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/config"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/internal/database"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/internal/routes"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize DynamoDB
	if err := database.InitDynamoDB(cfg); err != nil {
		log.Fatalf("Failed to initialize DynamoDB: %v", err)
	}

	// Create Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router)

	// Start server
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
