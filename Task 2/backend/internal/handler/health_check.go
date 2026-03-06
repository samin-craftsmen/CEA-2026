package handlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func HealthCheck(_ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{"status":"healthy"}`,
	}, nil
}
