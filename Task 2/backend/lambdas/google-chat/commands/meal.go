package commands

import (
	"fmt"
	"strings"

	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/client"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/types"
)

func init() {
	Register("meal", HandleMeal)
}

func HandleMeal(args string, userID string) *types.Response {
	userID, errResp := validatedCommandUserID(userID)
	if errResp != nil {
		return errResp
	}
	parts := strings.Fields(args)
	if len(parts) == 0 {
		return types.TextResponse("Usage:\nmeal view <date>\nmeal set <date> <meal_type> <YES|NO>")
	}

	switch normalizedSubcommand(parts[0]) {
	case "view":
		if resp := requireExactArgs(parts, 2, "meal view YYYY-MM-DD"); resp != nil {
			return resp
		}
		date, errResp := validatedChatDate(parts[1])
		if errResp != nil {
			return errResp
		}
		return handleMealView(userID, date)
	case "set":
		if resp := requireExactArgs(parts, 4, "meal set <date> <meal_type> <YES|NO>"); resp != nil {
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
		status, errResp := validatedChatStatus(parts[3])
		if errResp != nil {
			return errResp
		}
		return handleMealSet(userID, date, mealType, status)
	default:
		return types.ErrorResponse(fmt.Sprintf("unknown subcommand %q", parts[0]))
	}
}

type mealViewRequest struct {
	DiscordID string `json:"discord_id"`
	Date      string `json:"date"`
}

type mealViewResponse struct {
	Date   string            `json:"date"`
	Meals  map[string]string `json:"meals"`
	UserID string            `json:"user_id"`
}

func handleMealView(userID, date string) *types.Response {
	var result mealViewResponse
	if err := client.Post("/meal/participation/view", mealViewRequest{DiscordID: userID, Date: date}, &result); err != nil {
		return types.ErrorResponse("failed to fetch meal participation: " + err.Error())
	}

	lines := []string{"Meal Participation - " + result.Date}
	for _, mealType := range sortedKeys(result.Meals) {
		lines = append(lines, fmt.Sprintf("%s: %s", capitalize(mealType), result.Meals[mealType]))
	}
	lines = append(lines, "User: "+result.UserID)
	return types.TextResponse(strings.Join(lines, "\n"))
}

type mealSetRequest struct {
	DiscordID string `json:"discord_id"`
	Date      string `json:"date"`
	MealType  string `json:"meal_type"`
	Status    string `json:"status"`
}

func handleMealSet(userID, date, mealType, status string) *types.Response {
	if err := client.Post("/meal/participation/set", mealSetRequest{
		DiscordID: userID,
		Date:      date,
		MealType:  mealType,
		Status:    strings.ToUpper(status),
	}, nil); err != nil {
		return types.ErrorResponse("failed to update meal status: " + err.Error())
	}

	verb := "opted in to"
	if strings.EqualFold(status, "NO") {
		verb = "opted out of"
	}
	return types.TextResponse(fmt.Sprintf("You have %s %s on %s.", verb, capitalize(strings.ToLower(mealType)), date))
}
