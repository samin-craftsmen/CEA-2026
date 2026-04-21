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

func HandleWorkLocation(args string, userID string) *types.Response {
	userID, errResp := validatedCommandUserID(userID)
	if errResp != nil {
		return errResp
	}
	parts := strings.Fields(args)
	if len(parts) == 0 {
		return types.TextResponse("Usage:\nwork-location view <date>\nwork-location set <date> <OFFICE|WFH>")
	}

	switch normalizedSubcommand(parts[0]) {
	case "view":
		if resp := requireExactArgs(parts, 2, "work-location view YYYY-MM-DD"); resp != nil {
			return resp
		}
		date, errResp := validatedChatDate(parts[1])
		if errResp != nil {
			return errResp
		}
		return handleWorkLocationView(userID, date)
	case "set":
		if resp := requireExactArgs(parts, 3, "work-location set <date> <OFFICE|WFH>"); resp != nil {
			return resp
		}
		date, errResp := validatedChatDate(parts[1])
		if errResp != nil {
			return errResp
		}
		location, errResp := validatedChatLocation(parts[2])
		if errResp != nil {
			return errResp
		}
		return handleWorkLocationSet(userID, date, location)
	default:
		return types.ErrorResponse(fmt.Sprintf("unknown subcommand %q", parts[0]))
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

	label := "Office"
	if strings.EqualFold(result.Location, "WFH") {
		label = "Work from Home"
	}
	return types.TextResponse(fmt.Sprintf("Work Location - %s\nLocation: %s", result.Date, label))
}

type workLocationSetRequest struct {
	DiscordID string `json:"discord_id"`
	Date      string `json:"date"`
	Location  string `json:"location"`
}

func handleWorkLocationSet(userID, date, location string) *types.Response {
	if err := client.Post("/work-location/set", workLocationSetRequest{DiscordID: userID, Date: date, Location: strings.ToUpper(location)}, nil); err != nil {
		return types.ErrorResponse("failed to update work location: " + err.Error())
	}

	label := "Office"
	mealNote := "Your meals have been opted in."
	if strings.EqualFold(location, "WFH") {
		label = "Work from Home"
		mealNote = "Your meals have been opted out for this date."
	}
	return types.TextResponse(fmt.Sprintf("Your work location on %s has been set to %s. %s", date, label, mealNote))
}
