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

func HandleMealType(args string, userID string) *types.Response {
	userID, errResp := validatedCommandUserID(userID)
	if errResp != nil {
		return errResp
	}
	parts := strings.Fields(args)
	if len(parts) == 0 {
		return types.TextResponse("Usage:\nmeal-type view <date>\nmeal-type add <date> <meal_type>")
	}

	switch normalizedSubcommand(parts[0]) {
	case "view":
		if resp := requireExactArgs(parts, 2, "meal-type view YYYY-MM-DD"); resp != nil {
			return resp
		}
		date, errResp := validatedChatDate(parts[1])
		if errResp != nil {
			return errResp
		}
		return handleMealTypeView(date)
	case "add":
		if resp := requireExactArgs(parts, 3, "meal-type add <date> <meal_type>"); resp != nil {
			return resp
		}
		date, errResp := validatedChatDate(parts[1])
		if errResp != nil {
			return errResp
		}
		mealType, errResp := validatedChatMealType(parts[2])
		if errResp != nil {
			return errResp
		}
		return handleMealTypeAdd(userID, date, mealType)
	default:
		return types.ErrorResponse(fmt.Sprintf("unknown subcommand %q", parts[0]))
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

	lines := []string{"Meal Types - " + date}
	for _, mealType := range result.MealTypes {
		lines = append(lines, "- "+capitalize(mealType))
	}
	return types.TextResponse(strings.Join(lines, "\n"))
}

type mealTypeAddRequest struct {
	AdminDiscordID string `json:"admin_discord_id"`
	Date           string `json:"date"`
	MealType       string `json:"meal_type"`
}

func handleMealTypeAdd(userID, date, mealType string) *types.Response {
	if err := client.Post("/meal/types/set", mealTypeAddRequest{AdminDiscordID: userID, Date: date, MealType: mealType}, nil); err != nil {
		return types.ErrorResponse("failed to add meal type: " + err.Error())
	}
	return types.TextResponse(fmt.Sprintf("%s has been added as a meal type for %s.", capitalize(strings.ToLower(mealType)), date))
}
