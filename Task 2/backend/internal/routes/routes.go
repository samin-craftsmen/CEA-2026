package routes

import (
	"github.com/aws/aws-lambda-go/events"
)

// SetupRoutes configures all API routes
func HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{}, nil
}
