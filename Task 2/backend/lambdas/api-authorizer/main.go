package main

import (
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const defaultAPIKey = "api-key-test"

func handler(request events.APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	providedKey := strings.TrimSpace(request.Headers["x-api-key"])
	if providedKey == "" {
		providedKey = strings.TrimSpace(request.Headers["X-API-Key"])
	}

	expectedKey := strings.TrimSpace(os.Getenv("CUSTOM_API_KEY"))
	if expectedKey == "" {
		expectedKey = defaultAPIKey
	}

	isAuthorized := providedKey != "" && providedKey == expectedKey

	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: isAuthorized,
		Context: map[string]interface{}{
			"authType": "api-key",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}