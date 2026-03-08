package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/common/config"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/common/database"
)

// HealthHandler returns a simple OK response
func HealthHandler() (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       "OK",
	}, nil
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize DynamoDB
	if err := database.InitDynamoDB(cfg); err != nil {
		log.Fatalf("Failed to initialize DynamoDB: %v", err)
	}

	// Start Lambda handler
	lambda.Start(HealthHandler)
}
