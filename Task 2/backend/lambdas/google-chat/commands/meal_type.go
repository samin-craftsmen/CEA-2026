package commands

import (
	"fmt"
	"strings"

	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/client"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/types"
)

func init() {
	Register("meal-type", HandleMealType)
}

// HandleMealType handles: meal-type view <date> | meal-type add <date> <meal_type>
func HandleMealType(args string, userID string) *types.Response {
	parts := strings.Fields(args)
	if len(parts) == 0 {
		return types.TextResponse("Usage:\n`/meal-type view <date>`\n`/meal-type add <date> <meal_type>`")
	}

	switch parts[0] {
	case "view":
		if len(parts) < 2 {
			return types.ErrorResponse("provide a date. Usage: `/meal-type view YYYY-MM-DD`")
		}
		return handleMealTypeView(parts[1])
	case "add":
		if len(parts) < 3 {
			return types.ErrorResponse("Usage: `/meal-type add <date> <meal_type>`")
		}
		return handleMealTypeAdd(userID, parts[1], parts[2])
	default:
		return types.ErrorResponse(fmt.Sprintf("unknown subcommand `%s`", parts[0]))
	}
}

type mealTypeViewRequest struct {
	Date string `json:"date"`
}

type mealTypeViewResponse struct {
	Date      string   `json:"date"`
	MealTypes []string `json:"meal_types"`
}

func handleMealTypeView(date string) *types.Response {
	var result mealTypeViewResponse
	if err := client.Post("/meal/types/view", mealTypeViewRequest{Date: date}, &result); err != nil {
		return types.ErrorResponse("failed to fetch meal types: " + err.Error())
	}

	if len(result.MealTypes) == 0 {
		return types.TextResponse("No meal types configured for " + date)
	}

	parts := make([]string, len(result.MealTypes))
	for i, mt := range result.MealTypes {
		parts[i] = "• " + capitalize(mt)
	}
	return types.CardResponse("Meal Types — "+date, "", []types.Widget{
		types.Para(strings.Join(parts, "\n")),
	})
}

type mealTypeAddRequest struct {
	AdminDiscordID string `json:"admin_discord_id"`
	Date           string `json:"date"`
	MealType       string `json:"meal_type"`
}

func handleMealTypeAdd(userID, date, mealType string) *types.Response {
	if err := client.Post("/meal/types/set", mealTypeAddRequest{
		AdminDiscordID: userID, Date: date, MealType: mealType,
	}, nil); err != nil {
		return types.ErrorResponse("failed to add meal type: " + err.Error())
	}
	return types.TextResponse(fmt.Sprintf("✅ *%s* has been added as a meal type for %s.", capitalize(strings.ToLower(mealType)), date))
}
