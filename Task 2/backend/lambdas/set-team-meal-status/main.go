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

type SetTeamMealStatusRequest struct {
	TeamLeadDiscordID string `json:"team_lead_discord_id"`
	TargetDiscordID   string `json:"target_discord_id"`
	Date              string `json:"date"`
	MealType          string `json:"meal_type"`
	Status            string `json:"status"`
}

func SetTeamMealStatusHandler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var req SetTeamMealStatusRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil ||
		req.TeamLeadDiscordID == "" || req.TargetDiscordID == "" ||
		req.Date == "" || req.MealType == "" || req.Status == "" {
		return jsonResponse(http.StatusBadRequest, map[string]string{
			"error": "invalid request: team_lead_discord_id, target_discord_id, date, meal_type, and status are required",
		}), nil
	}

	if err := service.TeamLeadSetMealStatus(req.TeamLeadDiscordID, req.TargetDiscordID, req.Date, req.MealType, req.Status); err != nil {
		var ve *service.ValidationError
		if errors.As(err, &ve) {
			return jsonResponse(http.StatusBadRequest, map[string]string{"error": err.Error()}), nil
		}
		return jsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}), nil
	}

	return jsonResponse(http.StatusOK, map[string]string{"message": "meal status updated successfully"}), nil
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
	lambda.Start(SetTeamMealStatusHandler)
}
