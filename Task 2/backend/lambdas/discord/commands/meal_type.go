package commands

import (
	"fmt"
	"strings"

	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/client"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/types"
)

func init() {
	Register("meal-type", HandleMealType)
}

// HandleMealType routes /meal-type subcommands.
func HandleMealType(data *types.CommandData, userID string) *types.InteractionResponse {
	if _, errResp := commandUserID(userID); errResp != nil {
		return errResp
	}
	if len(data.Options) == 0 {
		return types.ErrorResponse("Please provide a subcommand. Usage: `/meal-type view` or `/meal-type add`")
	}

	sub := data.Options[0]
	switch sub.Name {
	case "view":
		date, errResp := validatedDateOption(sub.Options, "date")
		if errResp != nil {
			return errResp
		}
		return handleMealTypeView(date)
	case "add":
		date, errResp := validatedDateOption(sub.Options, "date")
		if errResp != nil {
			return errResp
		}
		mealType, errResp := validatedMealTypeOption(sub.Options, "meal_type")
		if errResp != nil {
			return errResp
		}
		return handleMealTypeAdd(userID, date, mealType)
	default:
		return types.ErrorResponse(fmt.Sprintf("Unknown subcommand: `%s`", sub.Name))
	}
}

type mealTypeViewRequest struct {
	Date string `json:"date"`
}

type mealTypeViewResponse struct {
	Date      string   `json:"date"`
	MealTypes []string `json:"meal_types"`
}

func handleMealTypeView(date string) *types.InteractionResponse {
	var result mealTypeViewResponse
	err := client.Post("/meal/types/view", mealTypeViewRequest{Date: date}, &result)
	if err != nil {
		return types.ErrorResponse("Failed to fetch meal types: " + err.Error())
	}

	if len(result.MealTypes) == 0 {
		return types.EmbedResponse(types.Embed{
			Title:       "Meal Types — " + date,
			Description: "No meal types configured.",
			Color:       types.ColorInfo,
		})
	}

	parts := make([]string, len(result.MealTypes))
	for i, mt := range result.MealTypes {
		parts[i] = fmt.Sprintf("• %s", capitalize(mt))
	}

	return types.EmbedResponse(types.Embed{
		Title:       "Meal Types — " + date,
		Description: strings.Join(parts, "\n"),
		Color:       types.ColorInfo,
	})
}

type mealTypeAddRequest struct {
	AdminDiscordID string `json:"admin_discord_id"`
	Date           string `json:"date"`
	MealType       string `json:"meal_type"`
}

func handleMealTypeAdd(userID, date, mealType string) *types.InteractionResponse {
	err := client.Post("/meal/types/set", mealTypeAddRequest{
		AdminDiscordID: userID,
		Date:           date,
		MealType:       mealType,
	}, nil)
	if err != nil {
		return types.ErrorResponse("Failed to add meal type: " + err.Error())
	}

	return types.EmbedResponse(types.Embed{
		Title:       "Meal Type Added",
		Description: fmt.Sprintf("✅ **%s** has been added as a meal type for %s.", capitalize(strings.ToLower(mealType)), date),
		Color:       types.ColorSuccess,
		Footer:      &types.EmbedFooter{Text: "Added by: " + userID},
	})
}
