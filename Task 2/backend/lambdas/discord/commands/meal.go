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
	if _, errResp := commandUserID(userID); errResp != nil {
		return errResp
	}
	if len(data.Options) == 0 {
		return types.ErrorResponse("Please provide a subcommand. Usage: `/meal view`")
	}

	sub := data.Options[0]
	switch sub.Name {
	case "view":
		date, errResp := validatedDateOption(sub.Options, "date")
		if errResp != nil {
			return errResp
		}
		return handleMealView(userID, date)
	case "set":
		date, errResp := validatedDateOption(sub.Options, "date")
		if errResp != nil {
			return errResp
		}
		mealType, errResp := validatedMealTypeOption(sub.Options, "meal_type")
		if errResp != nil {
			return errResp
		}
		status, errResp := validatedStatusOption(sub.Options, "status")
		if errResp != nil {
			return errResp
		}
		return handleMealSet(userID, date, mealType, status)
	default:
		return types.ErrorResponse(fmt.Sprintf("Unknown subcommand: `%s`", sub.Name))
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

func handleMealView(userID string, date string) *types.InteractionResponse {
	var result mealViewResponse
	err := client.Post("/meal/participation/view", mealViewRequest{DiscordID: userID, Date: date}, &result)
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
		Title:  "Meal Participation — " + result.Date,
		Color:  types.ColorSuccess,
		Fields: fields,
		Footer: &types.EmbedFooter{Text: "User: " + result.UserID},
	})
}

type mealSetRequest struct {
	DiscordID string `json:"discord_id"`
	Date      string `json:"date"`
	MealType  string `json:"meal_type"`
	Status    string `json:"status"`
}

func handleMealSet(userID, date, mealType, status string) *types.InteractionResponse {
	err := client.Post("/meal/participation/set", mealSetRequest{
		DiscordID: userID,
		Date:      date,
		MealType:  mealType,
		Status:    status,
	}, nil)
	if err != nil {
		return types.ErrorResponse("Failed to update meal status: " + err.Error())
	}

	emoji := "✅"
	optText := "opted **in** to"
	if strings.EqualFold(status, "NO") {
		emoji = "❌"
		optText = "opted **out** of"
	}

	return types.EmbedResponse(types.Embed{
		Title:       "Meal Status Updated",
		Description: fmt.Sprintf("%s You have %s **%s** on %s.", emoji, optText, capitalize(mealType), date),
		Color:       types.ColorSuccess,
		Footer:      &types.EmbedFooter{Text: "User: " + userID},
	})
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
