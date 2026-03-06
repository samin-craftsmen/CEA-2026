package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
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

	// Start Lambda handler
	lambda.Start(routes.HandleRequest)
}
