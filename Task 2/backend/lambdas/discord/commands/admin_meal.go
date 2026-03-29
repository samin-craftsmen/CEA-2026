package commands

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/client"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/types"
)

func init() {
	Register("admin-meal", HandleAdminMeal)
}

// extractAPIErrorMessage parses API error responses to get the actual error message.
// API errors come in format: "API /path returned status 400: {"error": "actual message"}"
func extractAPIErrorMessage(err error) string {
	errStr := err.Error()

	// Find the JSON part after the colon
	colonIndex := strings.LastIndex(errStr, ": ")
	if colonIndex == -1 {
		return errStr
	}

	jsonStr := errStr[colonIndex+2:]

	// Try to parse as JSON
	var errorResp struct {
		Error string `json:"error"`
	}
	if json.Unmarshal([]byte(jsonStr), &errorResp) == nil && errorResp.Error != "" {
		return errorResp.Error
	}

	// Fallback to original error if parsing fails
	return errStr
}

// validateDiscordUserID checks if the provided string is a valid Discord user ID (numeric snowflake).
func validateDiscordUserID(userID string) bool {
	if userID == "" {
		return false
	}
	matched, _ := regexp.MatchString(`^\d+$`, userID)
	return matched
}

// validateDateFormat checks if the provided string is in YYYY-MM-DD format.
func validateDateFormat(date string) bool {
	if date == "" {
		return false
	}
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

// validateMealType checks if the provided meal type is valid.
func validateMealType(mealType string) bool {
	return mealType == "lunch" || mealType == "snacks"
}

// validateMealStatus checks if the provided meal status is valid.
func validateMealStatus(status string) bool {
	return status == "YES" || status == "NO"
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
		if !validateDiscordUserID(targetUserID) {
			return types.ErrorResponse("Invalid employee ID format. Please provide a valid Discord user ID.")
		}
		if !validateDateFormat(date) {
			return types.ErrorResponse("Invalid date format. Please use YYYY-MM-DD format (e.g. 2026-03-15).")
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
		if !validateDiscordUserID(targetUserID) {
			return types.ErrorResponse("Invalid employee ID format. Please provide a valid Discord user ID.")
		}
		if !validateDateFormat(date) {
			return types.ErrorResponse("Invalid date format. Please use YYYY-MM-DD format (e.g. 2026-03-15).")
		}
		if !validateMealType(mealType) {
			return types.ErrorResponse("Invalid meal type. Must be 'lunch' or 'snacks'.")
		}
		if !validateMealStatus(status) {
			return types.ErrorResponse("Invalid status. Must be 'YES' or 'NO'.")
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
		return types.ErrorResponse("Failed to fetch meal participation: " + extractAPIErrorMessage(err))
	}

	fields := make([]types.EmbedField, 0, len(result.Meals))
	for _, mt := range []string{"lunch", "snacks"} {
		participation, ok := result.Meals[mt]
		if !ok {
			continue
		}
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
		return types.ErrorResponse("Failed to update meal status: " + extractAPIErrorMessage(err))
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
