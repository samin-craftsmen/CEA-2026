package commands

import (
	"fmt"

	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/client"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/types"
)

func init() {
	Register("work-location", HandleWorkLocation)
}

// HandleWorkLocation routes /work-location subcommands.
func HandleWorkLocation(data *types.CommandData, userID string) *types.InteractionResponse {
	if len(data.Options) == 0 {
		return types.ErrorResponse("Please provide a subcommand. Usage: `/work-location view` or `/work-location set`")
	}

	sub := data.Options[0]
	switch sub.Name {
	case "view":
		var date string
		for _, opt := range sub.Options {
			if v, ok := opt.Value.(string); ok && opt.Name == "date" {
				date = v
			}
		}
		if date == "" {
			return types.ErrorResponse("Please provide the required option: date.")
		}
		return handleWorkLocationView(userID, date)
	case "set":
		var date, location string
		for _, opt := range sub.Options {
			if v, ok := opt.Value.(string); ok {
				switch opt.Name {
				case "date":
					date = v
				case "location":
					location = v
				}
			}
		}
		if date == "" || location == "" {
			return types.ErrorResponse("Please provide all required options: date and location.")
		}
		return handleWorkLocationSet(userID, date, location)
	default:
		return types.ErrorResponse(fmt.Sprintf("Unknown subcommand: `%s`", sub.Name))
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

func handleWorkLocationView(userID, date string) *types.InteractionResponse {
	var result workLocationViewResponse
	err := client.Post("/work-location/view", workLocationViewRequest{
		DiscordID: userID,
		Date:      date,
	}, &result)
	if err != nil {
		return types.ErrorResponse("Failed to fetch work location: " + err.Error())
	}

	emoji := "🏢"
	label := "Office"
	if result.Location == "WFH" {
		emoji = "🏠"
		label = "Work from Home"
	}

	return types.EmbedResponse(types.Embed{
		Title:       "Work Location — " + result.Date,
		Description: fmt.Sprintf("%s **%s**", emoji, label),
		Color:       types.ColorInfo,
	})
}

type workLocationSetRequest struct {
	DiscordID string `json:"discord_id"`
	Date      string `json:"date"`
	Location  string `json:"location"`
}

func handleWorkLocationSet(userID, date, location string) *types.InteractionResponse {
	err := client.Post("/work-location/set", workLocationSetRequest{
		DiscordID: userID,
		Date:      date,
		Location:  location,
	}, nil)
	if err != nil {
		return types.ErrorResponse("Failed to update work location: " + err.Error())
	}

	emoji := "🏢"
	label := "Office"
	mealNote := "Your meals have been opted **in**."
	if location == "WFH" {
		emoji = "🏠"
		label = "Work from Home"
		mealNote = "Your meals have been opted **out** for this date."
	}

	return types.EmbedResponse(types.Embed{
		Title:       "Work Location Updated",
		Description: fmt.Sprintf("%s Your work location on **%s** has been set to **%s**.\n%s", emoji, date, label, mealNote),
		Color:       types.ColorSuccess,
	})
}
