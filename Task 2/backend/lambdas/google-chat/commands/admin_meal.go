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

// HandleAdminMeal handles:
//
//	admin-meal view <target_user_id> <date>
//	admin-meal set <target_user_id> <date> <meal_type> <YES|NO>
//	admin-meal headcount <date>
func HandleAdminMeal(args string, userID string) *types.Response {
	parts := strings.Fields(args)
	if len(parts) == 0 {
		return types.TextResponse("Usage:\n`/admin-meal view <user_id> <date>`\n`/admin-meal set <user_id> <date> <meal_type> <YES|NO>`\n`/admin-meal headcount <date>`")
	}

	switch parts[0] {
	case "view":
		if len(parts) < 3 {
			return types.ErrorResponse("Usage: `/admin-meal view <user_id> <date>`")
		}
		return handleAdminMealView(userID, parts[1], parts[2])
	case "set":
		if len(parts) < 5 {
			return types.ErrorResponse("Usage: `/admin-meal set <user_id> <date> <meal_type> <YES|NO>`")
		}
		return handleAdminMealSet(userID, parts[1], parts[2], parts[3], parts[4])
	case "headcount":
		if len(parts) < 2 {
			return types.ErrorResponse("Usage: `/admin-meal headcount <date>`")
		}
		return handleAdminHeadcount(userID, parts[1])
	default:
		return types.ErrorResponse(fmt.Sprintf("unknown subcommand `%s`", parts[0]))
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
		AdminDiscordID: adminID, TargetDiscordID: targetUserID, Date: date,
	}, &result); err != nil {
		return types.ErrorResponse("failed to fetch meal participation: " + err.Error())
	}

	mealKeys := make([]string, 0, len(result.Meals))
	for mt := range result.Meals {
		mealKeys = append(mealKeys, mt)
	}
	sort.Strings(mealKeys)

	widgets := make([]types.Widget, 0, len(result.Meals))
	for _, mt := range mealKeys {
		status := result.Meals[mt]
		emoji := "✅"
		if strings.EqualFold(status, "NO") {
			emoji = "❌"
		}
		widgets = append(widgets, types.KV(capitalize(mt), emoji+" "+status))
	}

	return types.CardResponse("Meal Participation — "+result.Date, "Employee: "+result.UserID, widgets)
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
		AdminDiscordID: adminID, TargetDiscordID: targetUserID,
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
	if err := client.Post("/meal/admin/headcount", adminHeadcountRequest{
		AdminDiscordID: adminID, Date: date,
	}, &result); err != nil {
		return types.ErrorResponse("failed to fetch headcount: " + err.Error())
	}

	widgets := []types.Widget{
		types.KV("Work Location", fmt.Sprintf("🏢 Office: %d  |  🏠 WFH: %d", result.WorkLocation.Office, result.WorkLocation.WFH)),
	}

	mealKeys := make([]string, 0, len(result.Summary))
	for mt := range result.Summary {
		mealKeys = append(mealKeys, mt)
	}
	sort.Strings(mealKeys)

	for _, mt := range mealKeys {
		entry := result.Summary[mt]
		widgets = append(widgets, types.KV(capitalize(mt), fmt.Sprintf("✅ Yes: %d  |  ❌ No: %d", entry.Yes, entry.No)))
	}

	return types.CardResponse("Meal Headcount — "+date, "", widgets)
}
