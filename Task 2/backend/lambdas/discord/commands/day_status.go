package commands

import (
	"fmt"

	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/client"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/types"
)

func init() {
	Register("day-status", HandleDayStatus)
}

// HandleDayStatus routes /day-status subcommands.
func HandleDayStatus(data *types.CommandData, userID string) *types.InteractionResponse {
	if _, errResp := commandUserID(userID); errResp != nil {
		return errResp
	}
	if len(data.Options) == 0 {
		return types.ErrorResponse("Please provide a subcommand. Usage: `/day-status set` or `/day-status view`")
	}

	sub := data.Options[0]
	switch sub.Name {
	case "set":
		date, errResp := validatedDateOption(sub.Options, "date")
		if errResp != nil {
			return errResp
		}
		statusType, errResp := validatedDayStatusOption(sub.Options, "type")
		if errResp != nil {
			return errResp
		}
		note, errResp := validatedNoteOption(sub.Options, "note", statusType == "SPECIAL_EVENT")
		if errResp != nil {
			return errResp
		}
		return handleDayStatusSet(userID, date, statusType, note)
	case "view":
		date, errResp := validatedDateOption(sub.Options, "date")
		if errResp != nil {
			return errResp
		}
		return handleDayStatusView(date)
	default:
		return types.ErrorResponse(fmt.Sprintf("Unknown subcommand: `%s`", sub.Name))
	}
}

type dayStatusSetRequest struct {
	AdminDiscordID string `json:"admin_discord_id"`
	Date           string `json:"date"`
	StatusType     string `json:"status_type"`
	Note           string `json:"note"`
}

func handleDayStatusSet(adminID, date, statusType, note string) *types.InteractionResponse {
	err := client.Post("/day/status/set", dayStatusSetRequest{
		AdminDiscordID: adminID,
		Date:           date,
		StatusType:     statusType,
		Note:           note,
	}, nil)
	if err != nil {
		return types.ErrorResponse("Failed to set day status: " + err.Error())
	}

	var emoji, label string
	switch statusType {
	case "GOVERNMENT_HOLIDAY":
		emoji = "🏛️"
		label = "Government Holiday"
	case "OFFICE_CLOSED":
		emoji = "🔒"
		label = "Office Closed"
	case "SPECIAL_EVENT":
		emoji = "🎉"
		label = "Special Event"
	default:
		emoji = "📅"
		label = statusType
	}

	description := fmt.Sprintf("%s **%s** has been marked as **%s**.", emoji, date, label)
	if statusType == "OFFICE_CLOSED" {
		description += "\nAll meals for this day have been automatically opted out for all employees."
	}
	if note != "" {
		description += fmt.Sprintf("\n📝 Note: %s", note)
	}

	return types.EmbedResponse(types.Embed{
		Title:       "Day Status Updated",
		Description: description,
		Color:       types.ColorSuccess,
		Footer:      &types.EmbedFooter{Text: "Set by admin: " + adminID},
	})
}

type dayStatusViewRequest struct {
	Date string `json:"date"`
}

type dayStatusViewResponse struct {
	Date  string `json:"date"`
	Type  string `json:"type"`
	Note  string `json:"note"`
	SetBy string `json:"set_by"`
}

func handleDayStatusView(date string) *types.InteractionResponse {
	var result dayStatusViewResponse
	err := client.Post("/day/status/view", dayStatusViewRequest{Date: date}, &result)
	if err != nil {
		return types.ErrorResponse("Failed to fetch day status: " + err.Error())
	}

	var emoji, label string
	switch result.Type {
	case "GOVERNMENT_HOLIDAY":
		emoji = "🏛️"
		label = "Government Holiday"
	case "OFFICE_CLOSED":
		emoji = "🔒"
		label = "Office Closed"
	case "SPECIAL_EVENT":
		emoji = "🎉"
		label = "Special Event"
	default:
		emoji = "✅"
		label = "Normal Working Day"
	}

	fields := []types.EmbedField{
		{Name: "Status", Value: fmt.Sprintf("%s %s", emoji, label), Inline: true},
	}
	if result.Note != "" {
		fields = append(fields, types.EmbedField{Name: "Note", Value: result.Note, Inline: false})
	}
	if result.SetBy != "" {
		fields = append(fields, types.EmbedField{Name: "Set by", Value: "<@" + result.SetBy + ">", Inline: true})
	}

	color := types.ColorInfo
	if result.Type == "OFFICE_CLOSED" || result.Type == "GOVERNMENT_HOLIDAY" {
		color = types.ColorError
	}

	return types.EmbedResponse(types.Embed{
		Title:  "Day Status — " + result.Date,
		Color:  color,
		Fields: fields,
	})
}
