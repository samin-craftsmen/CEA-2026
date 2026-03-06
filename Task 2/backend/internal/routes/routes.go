package routes

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	handlers "github.com/samin-craftsmen/meal-headcount-planner-backend/internal/handler"
)

// SetupRoutes configures all API routes
func HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.HTTPMethod == http.MethodGet && req.Path == "/health" {
		return handlers.HealthCheck(req)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{"message":"not found"}`,
	}, nil
}
