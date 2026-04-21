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

func HandleTeamMeal(args string, userID string) *types.Response {
	userID, errResp := validatedCommandUserID(userID)
	if errResp != nil {
		return errResp
	}
	parts := strings.Fields(args)
	if len(parts) == 0 {
		return types.TextResponse("Usage:\nteam-meal view <date>\nteam-meal set <user_id> <date> <meal_type> <YES|NO>")
	}

	switch normalizedSubcommand(parts[0]) {
	case "view":
		if resp := requireExactArgs(parts, 2, "team-meal view YYYY-MM-DD"); resp != nil {
			return resp
		}
		date, errResp := validatedChatDate(parts[1])
		if errResp != nil {
			return errResp
		}
		return handleTeamMealView(userID, date)
	case "set":
		if resp := requireExactArgs(parts, 5, "team-meal set <user_id> <date> <meal_type> <YES|NO>"); resp != nil {
			return resp
		}
		targetUserID, errResp := validatedChatTargetUserID(parts[1])
		if errResp != nil {
			return errResp
		}
		date, errResp := validatedChatDate(parts[2])
		if errResp != nil {
			return errResp
		}
		mealType, errResp := validatedChatMealType(parts[3])
		if errResp != nil {
			return errResp
		}
		status, errResp := validatedChatStatus(parts[4])
		if errResp != nil {
			return errResp
		}
		return handleTeamMealSet(userID, targetUserID, date, mealType, status)
	default:
		return types.ErrorResponse(fmt.Sprintf("unknown subcommand %q", parts[0]))
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
	for memberID := range result.Members {
		memberIDs = append(memberIDs, memberID)
	}
	sort.Strings(memberIDs)

	lines := []string{"Team Meal Status - " + date}
	if result.TeamID != "" {
		lines = append(lines, "Team: "+result.TeamID)
	}
	for _, memberID := range memberIDs {
		mealStates := result.Members[memberID]
		mealKeys := make([]string, 0, len(mealStates))
		for mealType := range mealStates {
			mealKeys = append(mealKeys, mealType)
		}
		sort.Strings(mealKeys)

		parts := make([]string, 0, len(mealKeys))
		for _, mealType := range mealKeys {
			parts = append(parts, fmt.Sprintf("%s: %s", capitalize(mealType), mealStates[mealType]))
		}
		lines = append(lines, fmt.Sprintf("%s -> %s", memberID, strings.Join(parts, " | ")))
	}

	return types.TextResponse(strings.Join(lines, "\n"))
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
		TeamLeadDiscordID: teamLeadID,
		TargetDiscordID:   targetUserID,
		Date:              date,
		MealType:          mealType,
		Status:            strings.ToUpper(status),
	}, nil); err != nil {
		return types.ErrorResponse("failed to update meal status: " + err.Error())
	}

	verb := "opted in to"
	if strings.EqualFold(status, "NO") {
		verb = "opted out of"
	}
	return types.TextResponse(fmt.Sprintf("User %s has been %s %s on %s.", targetUserID, verb, capitalize(strings.ToLower(mealType)), date))
}
