package commands

import (
	"fmt"
	"strings"

	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/client"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/types"
)

func init() {
	Register("work-location", HandleWorkLocation)
}

// HandleWorkLocation handles: work-location view <date> | work-location set <date> <OFFICE|WFH>
func HandleWorkLocation(args string, userID string) *types.Response {
	parts := strings.Fields(args)
	if len(parts) == 0 {
		return types.TextResponse("Usage:\n`/work-location view <date>`\n`/work-location set <date> <OFFICE|WFH>`")
	}

	switch parts[0] {
	case "view":
		if len(parts) < 2 {
			return types.ErrorResponse("provide a date. Usage: `/work-location view YYYY-MM-DD`")
		}
		return handleWorkLocationView(userID, parts[1])
	case "set":
		if len(parts) < 3 {
			return types.ErrorResponse("Usage: `/work-location set <date> <OFFICE|WFH>`")
		}
		return handleWorkLocationSet(userID, parts[1], parts[2])
	default:
		return types.ErrorResponse(fmt.Sprintf("unknown subcommand `%s`", parts[0]))
	}
}

type workLocationViewRequest struct {
	DiscordID string `json:"discord_id"`
	Date      string `json:"date"`
}

type workLocationViewResponse struct {
	Date     string `json:"date"`
	UserID   string `json:"user_id"`
	Location string `json:"location"`
}

func handleWorkLocationView(userID, date string) *types.Response {
	var result workLocationViewResponse
	if err := client.Post("/work-location/view", workLocationViewRequest{DiscordID: userID, Date: date}, &result); err != nil {
		return types.ErrorResponse("failed to fetch work location: " + err.Error())
	}

	emoji, label := "🏢", "Office"
	if result.Location == "WFH" {
		emoji, label = "🏠", "Work from Home"
	}
	return types.CardResponse("Work Location — "+result.Date, "", []types.Widget{
		types.KV("Location", emoji+" "+label),
	})
}

type workLocationSetRequest struct {
	DiscordID string `json:"discord_id"`
	Date      string `json:"date"`
	Location  string `json:"location"`
}

func handleWorkLocationSet(userID, date, location string) *types.Response {
	if err := client.Post("/work-location/set", workLocationSetRequest{
		DiscordID: userID, Date: date, Location: location,
	}, nil); err != nil {
		return types.ErrorResponse("failed to update work location: " + err.Error())
	}

	emoji, label, mealNote := "🏢", "Office", "Your meals have been opted *in*."
	if strings.EqualFold(location, "WFH") {
		emoji, label, mealNote = "🏠", "Work from Home", "Your meals have been opted *out* for this date."
	}
	return types.TextResponse(fmt.Sprintf("%s Your work location on *%s* has been set to *%s*.\n%s", emoji, date, label, mealNote))
}
