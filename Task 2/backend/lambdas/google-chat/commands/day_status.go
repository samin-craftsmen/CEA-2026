package commands

import (
	"fmt"
	"strings"

	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/client"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/types"
)

func init() {
	Register("day-status", HandleDayStatus)
}

// HandleDayStatus handles:
//
//	day-status set <date> <GOVERNMENT_HOLIDAY|OFFICE_CLOSED|SPECIAL_EVENT> [note]
//	day-status view <date>
func HandleDayStatus(args string, userID string) *types.Response {
	parts := strings.Fields(args)
	if len(parts) == 0 {
		return types.TextResponse("Usage:\n`/day-status set <date> <type> [note]`\n`/day-status view <date>`")
	}

	switch parts[0] {
	case "set":
		if len(parts) < 3 {
			return types.ErrorResponse("Usage: `/day-status set <date> <GOVERNMENT_HOLIDAY|OFFICE_CLOSED|SPECIAL_EVENT> [note]`")
		}
		note := strings.Join(parts[3:], " ")
		return handleDayStatusSet(userID, parts[1], parts[2], note)
	case "view":
		if len(parts) < 2 {
			return types.ErrorResponse("provide a date. Usage: `/day-status view YYYY-MM-DD`")
		}
		return handleDayStatusView(parts[1])
	default:
		return types.ErrorResponse(fmt.Sprintf("unknown subcommand `%s`", parts[0]))
	}
}

type dayStatusSetRequest struct {
	AdminDiscordID string `json:"admin_discord_id"`
	Date           string `json:"date"`
	StatusType     string `json:"status_type"`
	Note           string `json:"note"`
}

func handleDayStatusSet(adminID, date, statusType, note string) *types.Response {
	if err := client.Post("/day/status/set", dayStatusSetRequest{
		AdminDiscordID: adminID, Date: date, StatusType: statusType, Note: note,
	}, nil); err != nil {
		return types.ErrorResponse("failed to set day status: " + err.Error())
	}

	emoji, label := dayStatusLabel(statusType)
	msg := fmt.Sprintf("%s *%s* has been marked as *%s*.", emoji, date, label)
	if statusType == "OFFICE_CLOSED" {
		msg += "\nAll meals for this day have been automatically opted out."
	}
	if note != "" {
		msg += "\n📝 Note: " + note
	}
	return types.TextResponse(msg)
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

func handleDayStatusView(date string) *types.Response {
	var result dayStatusViewResponse
	if err := client.Post("/day/status/view", dayStatusViewRequest{Date: date}, &result); err != nil {
		return types.ErrorResponse("failed to fetch day status: " + err.Error())
	}

	emoji, label := dayStatusLabel(result.Type)
	widgets := []types.Widget{
		types.KV("Status", emoji+" "+label),
	}
	if result.Note != "" {
		widgets = append(widgets, types.KV("Note", result.Note))
	}
	if result.SetBy != "" {
		widgets = append(widgets, types.KV("Set by", result.SetBy))
	}
	return types.CardResponse("Day Status — "+result.Date, "", widgets)
}

func dayStatusLabel(statusType string) (emoji, label string) {
	switch statusType {
	case "GOVERNMENT_HOLIDAY":
		return "🏛️", "Government Holiday"
	case "OFFICE_CLOSED":
		return "🔒", "Office Closed"
	case "SPECIAL_EVENT":
		return "🎉", "Special Event"
	default:
		return "✅", "Normal Working Day"
	}
}
