package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/client"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/types"
)

func init() {
	Register("admin-meal", HandleAdminMeal)
}

// HandleAdminMeal routes /admin-meal subcommands.
func HandleAdminMeal(data *types.CommandData, userID string) *types.InteractionResponse {
	if _, errResp := commandUserID(userID); errResp != nil {
		return errResp
	}
	if len(data.Options) == 0 {
		return types.ErrorResponse("Please provide a subcommand. Usage: `/admin-meal view`, `/admin-meal set`, or `/admin-meal headcount`")
	}

	sub := data.Options[0]
	switch sub.Name {
	case "view":
		targetUserID, errResp := validatedTargetUserIDOption(sub.Options, "employee")
		if errResp != nil {
			return errResp
		}
		date, errResp := validatedDateOption(sub.Options, "date")
		if errResp != nil {
			return errResp
		}
		return handleAdminMealView(userID, targetUserID, date)
	case "set":
		targetUserID, errResp := validatedTargetUserIDOption(sub.Options, "employee")
		if errResp != nil {
			return errResp
		}
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
		return handleAdminMealSet(userID, targetUserID, date, mealType, status)
	case "headcount":
		date, errResp := validatedDateOption(sub.Options, "date")
		if errResp != nil {
			return errResp
		}
		return handleAdminHeadcount(userID, date)
	default:
		return types.ErrorResponse(fmt.Sprintf("Unknown subcommand: `%s`", sub.Name))
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

func handleAdminMealView(adminID, targetUserID, date string) *types.InteractionResponse {
	var result adminMealViewResponse
	err := client.Post("/meal/admin/view", adminMealViewRequest{
		AdminDiscordID:  adminID,
		TargetDiscordID: targetUserID,
		Date:            date,
	}, &result)
	if err != nil {
		return types.ErrorResponse("Failed to fetch meal participation: " + err.Error())
	}

	mealTypeKeys := make([]string, 0, len(result.Meals))
	for mt := range result.Meals {
		mealTypeKeys = append(mealTypeKeys, mt)
	}
	sort.Strings(mealTypeKeys)

	fields := make([]types.EmbedField, 0, len(result.Meals))
	for _, mt := range mealTypeKeys {
		participation := result.Meals[mt]
		emoji := "✅"
		if strings.EqualFold(participation, "NO") {
			emoji = "❌"
		}
		fields = append(fields, types.EmbedField{
			Name:   capitalize(mt),
			Value:  fmt.Sprintf("%s %s", emoji, participation),
			Inline: true,
		})
	}

	return types.EmbedResponse(types.Embed{
		Title:  "Meal Participation — " + result.Date,
		Color:  types.ColorInfo,
		Fields: fields,
		Footer: &types.EmbedFooter{Text: "Employee: " + result.UserID},
	})
}

type adminMealSetRequest struct {
	AdminDiscordID  string `json:"admin_discord_id"`
	TargetDiscordID string `json:"target_discord_id"`
	Date            string `json:"date"`
	MealType        string `json:"meal_type"`
	Status          string `json:"status"`
}

func handleAdminMealSet(adminID, targetUserID, date, mealType, status string) *types.InteractionResponse {
	err := client.Post("/meal/admin/set", adminMealSetRequest{
		AdminDiscordID:  adminID,
		TargetDiscordID: targetUserID,
		Date:            date,
		MealType:        mealType,
		Status:          status,
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
		Title:       "Admin Meal Status Updated",
		Description: fmt.Sprintf("%s <@%s> has been %s **%s** on %s.", emoji, targetUserID, optText, capitalize(mealType), date),
		Color:       types.ColorSuccess,
		Footer:      &types.EmbedFooter{Text: "Updated by admin: " + adminID},
	})
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

func handleAdminHeadcount(adminID, date string) *types.InteractionResponse {
	var result adminHeadcountResponse
	err := client.Post("/meal/admin/headcount", adminHeadcountRequest{
		AdminDiscordID: adminID,
		Date:           date,
	}, &result)
	if err != nil {
		return types.ErrorResponse("Failed to fetch headcount: " + err.Error())
	}

	fields := []types.EmbedField{
		{
			Name:   "Work Location",
			Value:  fmt.Sprintf("🏢 Office: **%d**  |  🏠 WFH: **%d**", result.WorkLocation.Office, result.WorkLocation.WFH),
			Inline: false,
		},
		{
			Name:   "​",
			Value:  "**Meal Participation**",
			Inline: false,
		},
	}

	mealTypeKeys := make([]string, 0, len(result.Summary))
	for mt := range result.Summary {
		mealTypeKeys = append(mealTypeKeys, mt)
	}
	sort.Strings(mealTypeKeys)

	for _, mt := range mealTypeKeys {
		entry := result.Summary[mt]
		fields = append(fields, types.EmbedField{
			Name:   capitalize(mt),
			Value:  fmt.Sprintf("✅ Yes: **%d**  |  ❌ No: **%d**", entry.Yes, entry.No),
			Inline: false,
		})
	}

	return types.EmbedResponse(types.Embed{
		Title:  "Meal Headcount — " + date,
		Color:  types.ColorInfo,
		Fields: fields,
	})
}
