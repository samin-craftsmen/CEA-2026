package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/client"
	ghandler "github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/handler"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/types"

	_ "github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/commands"
)

func handler(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("Handler triggered")
	log.Printf("request: method=%s path=%s content-type=%q user-agent=%q auth-present=%t", req.RequestContext.HTTP.Method, req.RawPath, req.Headers["content-type"], req.Headers["user-agent"], req.Headers["authorization"] != "")
	log.Printf("raw body: %s", req.Body)

	if strings.TrimSpace(req.Body) == "" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"text":"Error: request body is required"}`,
		}, nil
	}

	var event types.Event
	if err := json.Unmarshal([]byte(req.Body), &event); err != nil {
		log.Println("parse error:", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"text":"Error: invalid request body"}`,
		}, nil
	}

	log.Printf("eventType=%q message=%q user=%q hostChat=%t", event.InteractionType(), event.MessageText(), event.GetUserID(), event.HostAppIsChat())

	response := ghandler.RouteEvent(&event)
	respBody := buildResponseBody(&event, response)
	log.Println("Returning:", respBody)

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       respBody,
	}, nil
}

func main() {
	client.Init()
	lambda.Start(handler)
}

func buildResponseBody(event *types.Event, response *types.Response) string {
	type chatMessage struct {
		Text string `json:"text"`
	}

	type createMessageAction struct {
		Message chatMessage `json:"message"`
	}

	type chatDataAction struct {
		CreateMessageAction createMessageAction `json:"createMessageAction"`
	}

	type hostAppDataAction struct {
		ChatDataAction chatDataAction `json:"chatDataAction"`
	}

	type dataActionsResponse struct {
		HostAppDataAction hostAppDataAction `json:"hostAppDataAction"`
	}

	responseText := ""
	if response != nil {
		responseText = response.Text
	}

	var payload any = chatMessage{Text: responseText}
	if event != nil && event.HostAppIsChat() {
		payload = dataActionsResponse{
			HostAppDataAction: hostAppDataAction{
				ChatDataAction: chatDataAction{
					CreateMessageAction: createMessageAction{
						Message: chatMessage{Text: responseText},
					},
				},
			},
		}
	}

	resp, err := json.Marshal(payload)
	if err != nil {
		log.Println("response marshal error:", err)
		return `{"text":""}`
	}
	return string(resp)
}
