package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/client"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/types"
)

func init() {
	Register("team-meal", HandleTeamMeal)
}

// HandleTeamMeal handles: team-meal view <date> | team-meal set <target_user_id> <date> <meal_type> <YES|NO>
func HandleTeamMeal(args string, userID string) *types.Response {
	parts := strings.Fields(args)
	if len(parts) == 0 {
		return types.TextResponse("Usage:\n`/team-meal view <date>`\n`/team-meal set <user_id> <date> <meal_type> <YES|NO>`")
	}

	switch parts[0] {
	case "view":
		if len(parts) < 2 {
			return types.ErrorResponse("provide a date. Usage: `/team-meal view YYYY-MM-DD`")
		}
		return handleTeamMealView(userID, parts[1])
	case "set":
		if len(parts) < 5 {
			return types.ErrorResponse("Usage: `/team-meal set <user_id> <date> <meal_type> <YES|NO>`")
		}
		return handleTeamMealSet(userID, parts[1], parts[2], parts[3], parts[4])
	default:
		return types.ErrorResponse(fmt.Sprintf("unknown subcommand `%s`", parts[0]))
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

func handleTeamMealView(teamLeadID, date string) *types.Response {
	var result teamMealViewResponse
	if err := client.Post("/meal/team/view", teamMealViewRequest{TeamLeadDiscordID: teamLeadID, Date: date}, &result); err != nil {
		return types.ErrorResponse("failed to fetch team meal status: " + err.Error())
	}

	if len(result.Members) == 0 {
		return types.TextResponse("No team members found for " + date)
	}

	memberIDs := make([]string, 0, len(result.Members))
	for id := range result.Members {
		memberIDs = append(memberIDs, id)
	}
	sort.Strings(memberIDs)

	widgets := make([]types.Widget, 0, len(result.Members))
	for _, memberID := range memberIDs {
		meals := result.Members[memberID]
		mealKeys := make([]string, 0, len(meals))
		for mt := range meals {
			mealKeys = append(mealKeys, mt)
		}
		sort.Strings(mealKeys)

		parts := make([]string, 0, len(meals))
		for _, mt := range mealKeys {
			emoji := "✅"
			if strings.EqualFold(meals[mt], "NO") {
				emoji = "❌"
			}
			parts = append(parts, fmt.Sprintf("%s %s: %s", emoji, capitalize(mt), meals[mt]))
		}
		widgets = append(widgets, types.KV("User "+memberID, strings.Join(parts, "  |  ")))
	}

	return types.CardResponse("Team Meal Status — "+date, result.TeamID, widgets)
}

type teamMealSetRequest struct {
	TeamLeadDiscordID string `json:"team_lead_discord_id"`
	TargetDiscordID   string `json:"target_discord_id"`
	Date              string `json:"date"`
	MealType          string `json:"meal_type"`
	Status            string `json:"status"`
}

func handleTeamMealSet(teamLeadID, targetUserID, date, mealType, status string) *types.Response {
	if err := client.Post("/meal/team/set", teamMealSetRequest{
		TeamLeadDiscordID: teamLeadID, TargetDiscordID: targetUserID,
		Date: date, MealType: mealType, Status: status,
	}, nil); err != nil {
		return types.ErrorResponse("failed to update meal status: " + err.Error())
	}

	emoji, optText := "✅", "opted *in* to"
	if strings.EqualFold(status, "NO") {
		emoji, optText = "❌", "opted *out* of"
	}
	return types.TextResponse(fmt.Sprintf("%s User %s has been %s *%s* on %s.", emoji, targetUserID, optText, capitalize(mealType), date))
}
