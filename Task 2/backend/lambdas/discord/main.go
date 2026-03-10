package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/client"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/handler"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/types"

	// Blank import to trigger init() command registration.
	_ "github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/commands"
)

var publicKey ed25519.PublicKey

// HandleInteraction is the Lambda entry point for Discord interactions.
func HandleInteraction(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// Verify request signature from Discord
	signature := request.Headers["x-signature-ed25519"]
	timestamp := request.Headers["x-signature-timestamp"]

	if !verifySignature(publicKey, signature, timestamp, request.Body) {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "invalid request signature",
		}, nil
	}

	// Parse the interaction payload
	var interaction types.Interaction
	if err := json.Unmarshal([]byte(request.Body), &interaction); err != nil {
		return jsonResponse(http.StatusBadRequest, types.ErrorResponse("Invalid request body")), nil
	}

	// Route to the appropriate handler
	response := handler.RouteInteraction(&interaction)
	return jsonResponse(http.StatusOK, response), nil
}

func verifySignature(key ed25519.PublicKey, signature, timestamp, body string) bool {
	if signature == "" || timestamp == "" {
		return false
	}

	sig, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	msg := []byte(timestamp + body)
	return ed25519.Verify(key, msg, sig)
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

	keyHex := os.Getenv("DISCORD_PUBLIC_KEY")
	if keyHex == "" {
		panic("DISCORD_PUBLIC_KEY environment variable is required")
	}

	key, err := hex.DecodeString(keyHex)
	if err != nil {
		panic("invalid DISCORD_PUBLIC_KEY: " + err.Error())
	}
	publicKey = ed25519.PublicKey(key)

	lambda.Start(HandleInteraction)
}
