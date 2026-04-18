package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/client"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/types"
)

func init() {
	Register("admin-meal", HandleAdminMeal)
}

func HandleAdminMeal(args string, userID string) *types.Response {
	parts := strings.Fields(args)
	if len(parts) == 0 {
		return types.TextResponse("Usage:\nadmin-meal view <user_id> <date>\nadmin-meal set <user_id> <date> <meal_type> <YES|NO>\nadmin-meal headcount <date>")
	}

	switch parts[0] {
	case "view":
		if len(parts) < 3 {
			return types.ErrorResponse("usage: admin-meal view <user_id> <date>")
		}
		return handleAdminMealView(userID, normalizeUserRef(parts[1]), parts[2])
	case "set":
		if len(parts) < 5 {
			return types.ErrorResponse("usage: admin-meal set <user_id> <date> <meal_type> <YES|NO>")
		}
		return handleAdminMealSet(userID, normalizeUserRef(parts[1]), parts[2], parts[3], parts[4])
	case "headcount":
		if len(parts) < 2 {
			return types.ErrorResponse("usage: admin-meal headcount <date>")
		}
		return handleAdminHeadcount(userID, parts[1])
	default:
		return types.ErrorResponse(fmt.Sprintf("unknown subcommand %q", parts[0]))
	}
}

type adminMealViewRequest struct {
	AdminDiscordID  string `json:"admin_discord_id"`
	TargetDiscordID string `json:"target_discord_id"`
	Date            string `json:"date"`
}

type adminMealViewResponse struct {
	Date   string            `json:"date"`
	UserID string            `json:"user_id"`
	Meals  map[string]string `json:"meals"`
}

func handleAdminMealView(adminID, targetUserID, date string) *types.Response {
	var result adminMealViewResponse
	if err := client.Post("/meal/admin/view", adminMealViewRequest{
		AdminDiscordID:  adminID,
		TargetDiscordID: targetUserID,
		Date:            date,
	}, &result); err != nil {
		return types.ErrorResponse("failed to fetch meal participation: " + err.Error())
	}

	mealKeys := make([]string, 0, len(result.Meals))
	for mealType := range result.Meals {
		mealKeys = append(mealKeys, mealType)
	}
	sort.Strings(mealKeys)

	lines := []string{"Meal Participation - " + result.Date, "Employee: " + result.UserID}
	for _, mealType := range mealKeys {
		lines = append(lines, fmt.Sprintf("%s: %s", capitalize(mealType), result.Meals[mealType]))
	}
	return types.TextResponse(strings.Join(lines, "\n"))
}

type adminMealSetRequest struct {
	AdminDiscordID  string `json:"admin_discord_id"`
	TargetDiscordID string `json:"target_discord_id"`
	Date            string `json:"date"`
	MealType        string `json:"meal_type"`
	Status          string `json:"status"`
}

func handleAdminMealSet(adminID, targetUserID, date, mealType, status string) *types.Response {
	if err := client.Post("/meal/admin/set", adminMealSetRequest{
		AdminDiscordID:  adminID,
		TargetDiscordID: targetUserID,
		Date:            date,
		MealType:        mealType,
		Status:          strings.ToUpper(status),
	}, nil); err != nil {
		return types.ErrorResponse("failed to update meal status: " + err.Error())
	}

	verb := "opted in to"
	if strings.EqualFold(status, "NO") {
		verb = "opted out of"
	}
	return types.TextResponse(fmt.Sprintf("User %s has been %s %s on %s.", targetUserID, verb, capitalize(strings.ToLower(mealType)), date))
}

type adminHeadcountRequest struct {
	AdminDiscordID string `json:"admin_discord_id"`
	Date           string `json:"date"`
}

type headcountEntry struct {
	Yes int `json:"yes"`
	No  int `json:"no"`
}

type workLocationSummary struct {
	Office int `json:"office"`
	WFH    int `json:"wfh"`
}

type adminHeadcountResponse struct {
	Date         string                    `json:"date"`
	WorkLocation workLocationSummary       `json:"work_location"`
	Summary      map[string]headcountEntry `json:"summary"`
}

func handleAdminHeadcount(adminID, date string) *types.Response {
	var result adminHeadcountResponse
	if err := client.Post("/meal/admin/headcount", adminHeadcountRequest{AdminDiscordID: adminID, Date: date}, &result); err != nil {
		return types.ErrorResponse("failed to fetch headcount: " + err.Error())
	}

	lines := []string{
		"Meal Headcount - " + date,
		fmt.Sprintf("Work Location: Office %d | WFH %d", result.WorkLocation.Office, result.WorkLocation.WFH),
	}
	mealKeys := make([]string, 0, len(result.Summary))
	for mealType := range result.Summary {
		mealKeys = append(mealKeys, mealType)
	}
	sort.Strings(mealKeys)
	for _, mealType := range mealKeys {
		entry := result.Summary[mealType]
		lines = append(lines, fmt.Sprintf("%s: YES %d | NO %d", capitalize(mealType), entry.Yes, entry.No))
	}
	return types.TextResponse(strings.Join(lines, "\n"))
}
