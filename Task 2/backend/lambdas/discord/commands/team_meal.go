package commands

import (
	"fmt"

	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/client"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/types"
)

func init() {
	Register("team-meal", HandleTeamMeal)
}

// HandleTeamMeal routes /team-meal subcommands.
func HandleTeamMeal(data *types.CommandData, userID string) *types.InteractionResponse {
	if len(data.Options) == 0 {
		return types.ErrorResponse("Please provide a subcommand. Usage: `/team-meal set`")
	}

	sub := data.Options[0]
	switch sub.Name {
	case "set":
		var targetUserID, date, mealType, status string
		for _, opt := range sub.Options {
			switch opt.Name {
			case "employee":
				// Discord USER option returns the user's snowflake ID as a string value.
				if v, ok := opt.Value.(string); ok {
					targetUserID = v
				}
			case "date":
				if v, ok := opt.Value.(string); ok {
					date = v
				}
			case "meal_type":
				if v, ok := opt.Value.(string); ok {
					mealType = v
				}
			case "status":
				if v, ok := opt.Value.(string); ok {
					status = v
				}
			}
		}
		if targetUserID == "" || date == "" || mealType == "" || status == "" {
			return types.ErrorResponse("Please provide all required options: employee, date, meal_type, and status.")
		}
		return handleTeamMealSet(userID, targetUserID, date, mealType, status)
	default:
		return types.ErrorResponse(fmt.Sprintf("Unknown subcommand: `%s`", sub.Name))
	}
}

type teamMealSetRequest struct {
	TeamLeadDiscordID string `json:"team_lead_discord_id"`
	TargetDiscordID   string `json:"target_discord_id"`
	Date              string `json:"date"`
	MealType          string `json:"meal_type"`
	Status            string `json:"status"`
}

func handleTeamMealSet(teamLeadID, targetUserID, date, mealType, status string) *types.InteractionResponse {
	err := client.Post("/meal/team/set", teamMealSetRequest{
		TeamLeadDiscordID: teamLeadID,
		TargetDiscordID:   targetUserID,
		Date:              date,
		MealType:          mealType,
		Status:            status,
	}, nil)
	if err != nil {
		return types.ErrorResponse("Failed to update meal status: " + err.Error())
	}

	emoji := "✅"
	optText := "opted **in** to"
	if status == "NO" {
		emoji = "❌"
		optText = "opted **out** of"
	}

	return types.EmbedResponse(types.Embed{
		Title:       "Team Meal Status Updated",
		Description: fmt.Sprintf("%s <@%s> has been %s **%s** on %s.", emoji, targetUserID, optText, capitalize(mealType), date),
		Color:       types.ColorSuccess,
		Footer:      &types.EmbedFooter{Text: "Updated by: " + teamLeadID},
	})
}
