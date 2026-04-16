package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/client"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/handler"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/types"

	// Blank imports to trigger init() command registration.
	_ "github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/commands"
)

var bearerToken string

// HandleInteraction is the Lambda entry point for Google Chat interactions.
func HandleInteraction(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// Verify the Bearer token sent by Google Chat.
	auth := request.Headers["authorization"]
	if !strings.EqualFold(strings.TrimSpace(strings.TrimPrefix(auth, "Bearer ")), bearerToken) {
		return jsonResponse(http.StatusUnauthorized, types.TextResponse("unauthorized")), nil
	}

	var event types.Event
	if err := json.Unmarshal([]byte(request.Body), &event); err != nil {
		return jsonResponse(http.StatusBadRequest, types.ErrorResponse("invalid request body")), nil
	}

	resp := handler.RouteEvent(&event)
	if resp == nil {
		// Acknowledge silently (e.g. non-command messages).
		return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusNoContent}, nil
	}
	return jsonResponse(http.StatusOK, resp), nil
}

func jsonResponse(statusCode int, body any) events.APIGatewayV2HTTPResponse {
	b, _ := json.Marshal(body)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: statusCode,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(b),
	}
}

func main() {
	client.Init()

	bearerToken = os.Getenv("GOOGLE_CHAT_BEARER_TOKEN")
	if bearerToken == "" {
		panic("GOOGLE_CHAT_BEARER_TOKEN environment variable is required")
	}

	lambda.Start(HandleInteraction)
}
