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
	if len(data.Options) == 0 {
		return types.ErrorResponse("Please provide a subcommand. Usage: `/admin-meal view` or `/admin-meal set`")
	}

	sub := data.Options[0]
	switch sub.Name {
	case "view":
		var targetUserID, date string
		for _, opt := range sub.Options {
			if v, ok := opt.Value.(string); ok {
				switch opt.Name {
				case "employee":
					targetUserID = v
				case "date":
					date = v
				}
			}
		}
		if targetUserID == "" || date == "" {
			return types.ErrorResponse("Please provide all required options: employee and date.")
		}
		return handleAdminMealView(userID, targetUserID, date)
	case "set":
		var targetUserID, date, mealType, status string
		for _, opt := range sub.Options {
			if v, ok := opt.Value.(string); ok {
				switch opt.Name {
				case "employee":
					targetUserID = v
				case "date":
					date = v
				case "meal_type":
					mealType = v
				case "status":
					status = v
				}
			}
		}
		if targetUserID == "" || date == "" || mealType == "" || status == "" {
			return types.ErrorResponse("Please provide all required options: employee, date, meal_type, and status.")
		}
		return handleAdminMealSet(userID, targetUserID, date, mealType, status)
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
