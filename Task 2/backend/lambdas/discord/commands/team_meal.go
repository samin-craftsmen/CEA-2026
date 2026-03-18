package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/client"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/types"
)

func init() {
	Register("team-meal", HandleTeamMeal)
}

// HandleTeamMeal routes /team-meal subcommands.
func HandleTeamMeal(data *types.CommandData, userID string) *types.InteractionResponse {
	if len(data.Options) == 0 {
		return types.ErrorResponse("Please provide a subcommand. Usage: `/team-meal view` or `/team-meal set`")
	}

	sub := data.Options[0]
	switch sub.Name {
	case "view":
		date := ""
		if len(sub.Options) > 0 {
			if v, ok := sub.Options[0].Value.(string); ok {
				date = v
			}
		}
		if date == "" {
			return types.ErrorResponse("Please provide a date.")
		}
		return handleTeamMealView(userID, date)
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

type teamMealViewRequest struct {
	TeamLeadDiscordID string `json:"team_lead_discord_id"`
	Date              string `json:"date"`
}

type teamMealViewResponse struct {
	Date    string                       `json:"date"`
	TeamID  string                       `json:"team_id"`
	Members map[string]map[string]string `json:"members"`
}

func handleTeamMealView(teamLeadID, date string) *types.InteractionResponse {
	var result teamMealViewResponse
	err := client.Post("/meal/team/view", teamMealViewRequest{
		TeamLeadDiscordID: teamLeadID,
		Date:              date,
	}, &result)
	if err != nil {
		return types.ErrorResponse("Failed to fetch team meal status: " + err.Error())
	}

	if len(result.Members) == 0 {
		return types.EmbedResponse(types.Embed{
			Title:       "Team Meal Status — " + date,
			Description: "No team members found.",
			Color:       types.ColorInfo,
		})
	}

	// One field per member showing all meal statuses on one line.
	fields := make([]types.EmbedField, 0, len(result.Members))

	// Sort member IDs for consistent ordering.
	memberIDs := make([]string, 0, len(result.Members))
	for id := range result.Members {
		memberIDs = append(memberIDs, id)
	}
	sort.Strings(memberIDs)

	for _, memberID := range memberIDs {
		meals := result.Members[memberID]
		parts := make([]string, 0, len(meals))
		for _, mt := range []string{"lunch", "snacks"} {
			status, ok := meals[mt]
			if !ok {
				continue
			}
			emoji := "✅"
			if strings.EqualFold(status, "NO") {
				emoji = "❌"
			}
			parts = append(parts, fmt.Sprintf("%s %s: %s", emoji, capitalize(mt), status))
		}
		fields = append(fields, types.EmbedField{
			Name:  fmt.Sprintf("<@%s>", memberID),
			Value: strings.Join(parts, "  |  "),
		})
	}

	return types.EmbedResponse(types.Embed{
		Title:  "Team Meal Status — " + date,
		Color:  types.ColorInfo,
		Fields: fields,
		Footer: &types.EmbedFooter{Text: result.TeamID},
	})
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
