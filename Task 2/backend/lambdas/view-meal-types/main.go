package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/common/config"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/common/database"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/common/service"
)

type ViewMealTypesRequest struct {
	Date string `json:"date"`
}

type ViewMealTypesResponse struct {
	Date      string   `json:"date"`
	MealTypes []string `json:"meal_types"`
}

func ViewMealTypesHandler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var req ViewMealTypesRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil || req.Date == "" {
		return jsonResponse(http.StatusBadRequest, map[string]string{
			"error": "invalid request: date is required",
		}), nil
	}

	mealTypes, err := service.GetMealTypesForDate(req.Date)
	if err != nil {
		return jsonResponse(http.StatusInternalServerError, map[string]string{"error": err.Error()}), nil
	}

	return jsonResponse(http.StatusOK, ViewMealTypesResponse{
		Date:      req.Date,
		MealTypes: mealTypes,
	}), nil
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
	lambda.Start(ViewMealTypesHandler)
}
