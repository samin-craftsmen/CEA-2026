package commands

import (
	"fmt"
	"strings"

	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/client"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/types"
)

func init() {
	Register("meal", HandleMeal)
}

// HandleMeal routes /meal subcommands.
func HandleMeal(data *types.CommandData, userID string) *types.InteractionResponse {
	if len(data.Options) == 0 {
		return types.ErrorResponse("Please provide a subcommand. Usage: `/meal view`")
	}

	sub := data.Options[0].Name
	switch sub {
	case "view":
		return handleMealView(userID)
	default:
		return types.ErrorResponse(fmt.Sprintf("Unknown subcommand: `%s`", sub))
	}
}

type mealViewRequest struct {
	DiscordID string `json:"discord_id"`
}

type mealViewResponse struct {
	Date   string            `json:"date"`
	Meals  map[string]string `json:"meals"`
	UserID string            `json:"user_id"`
}

func handleMealView(userID string) *types.InteractionResponse {
	var result mealViewResponse
	err := client.Post("/meal/participation/view", mealViewRequest{DiscordID: userID}, &result)
	if err != nil {
		return types.ErrorResponse("Failed to fetch meal participation: " + err.Error())
	}

	fields := make([]types.EmbedField, 0, len(result.Meals))
	for mealType, participation := range result.Meals {
		emoji := "✅"
		if strings.EqualFold(participation, "NO") {
			emoji = "❌"
		}
		fields = append(fields, types.EmbedField{
			Name:   capitalize(mealType),
			Value:  fmt.Sprintf("%s %s", emoji, participation),
			Inline: true,
		})
	}

	return types.EmbedResponse(types.Embed{
		Title:  "🍽️ Meal Participation — " + result.Date,
		Color:  types.ColorSuccess,
		Fields: fields,
		Footer: &types.EmbedFooter{Text: "User: " + result.UserID},
	})
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
