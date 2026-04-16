package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/client"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/types"
)

func init() {
	Register("meal", HandleMeal)
}

// HandleMeal handles: meal view <date> | meal set <date> <meal_type> <YES|NO>
func HandleMeal(args string, userID string) *types.Response {
	parts := strings.Fields(args)
	if len(parts) == 0 {
		return types.TextResponse("Usage:\n`/meal view <date>`\n`/meal set <date> <meal_type> <YES|NO>`")
	}

	switch parts[0] {
	case "view":
		if len(parts) < 2 {
			return types.ErrorResponse("provide a date. Usage: `/meal view YYYY-MM-DD`")
		}
		return handleMealView(userID, parts[1])
	case "set":
		if len(parts) < 4 {
			return types.ErrorResponse("Usage: `/meal set <date> <meal_type> <YES|NO>`")
		}
		return handleMealSet(userID, parts[1], parts[2], parts[3])
	default:
		return types.ErrorResponse(fmt.Sprintf("unknown subcommand `%s`", parts[0]))
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

	mealKeys := sortedKeys(result.Meals)
	widgets := make([]types.Widget, 0, len(result.Meals))
	for _, mt := range mealKeys {
		status := result.Meals[mt]
		emoji := "✅"
		if strings.EqualFold(status, "NO") {
			emoji = "❌"
		}
		widgets = append(widgets, types.KV(capitalize(mt), emoji+" "+status))
	}

	return types.CardResponse("Meal Participation — "+result.Date, "User: "+result.UserID, widgets)
}

type mealSetRequest struct {
	DiscordID string `json:"discord_id"`
	Date      string `json:"date"`
	MealType  string `json:"meal_type"`
	Status    string `json:"status"`
}

func handleMealSet(userID, date, mealType, status string) *types.Response {
	if err := client.Post("/meal/participation/set", mealSetRequest{
		DiscordID: userID, Date: date, MealType: mealType, Status: status,
	}, nil); err != nil {
		return types.ErrorResponse("failed to update meal status: " + err.Error())
	}

	emoji, optText := "✅", "opted *in* to"
	if strings.EqualFold(status, "NO") {
		emoji, optText = "❌", "opted *out* of"
	}
	return types.TextResponse(fmt.Sprintf("%s You have %s *%s* on %s.", emoji, optText, capitalize(mealType), date))
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
