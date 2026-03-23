package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/common/config"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/common/database"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/common/service"
)

type SetDayStatusRequest struct {
	AdminDiscordID string `json:"admin_discord_id"`
	Date           string `json:"date"`
	StatusType     string `json:"status_type"`
	Note           string `json:"note"`
}

func SetDayStatusHandler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var req SetDayStatusRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil ||
		req.AdminDiscordID == "" || req.Date == "" || req.StatusType == "" {
		return jsonResponse(http.StatusBadRequest, map[string]string{
			"error": "invalid request: admin_discord_id, date, and status_type are required",
		}), nil
	}

	if err := service.AdminSetDayStatus(req.AdminDiscordID, req.Date, req.StatusType, req.Note); err != nil {
		var ve *service.ValidationError
		if errors.As(err, &ve) {
			return jsonResponse(http.StatusBadRequest, map[string]string{"error": err.Error()}), nil
		}
		return jsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}), nil
	}

	return jsonResponse(http.StatusOK, map[string]string{"message": "day status set successfully"}), nil
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
	cfg := config.LoadConfig()
	if err := database.InitDynamoDB(cfg); err != nil {
		panic(err)
	}
	lambda.Start(SetDayStatusHandler)
}
