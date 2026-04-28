package handler

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/commands"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/types"
)

var (
	commandIDsOnce sync.Once
	commandIDMap   map[int64]string
)

func RouteEvent(event *types.Event) *types.Response {
	if event == nil {
		return types.ErrorResponse("invalid event payload")
	}

	switch event.InteractionType() {
	case types.InteractionTypeAddedToSpace:
		return types.TextResponse(welcomeText())
	case types.InteractionTypeMessage:
		if strings.TrimSpace(event.GetUserID()) == "" {
			return types.ErrorResponse("unable to determine the invoking user")
		}
		return handleCommandText(event.MessageText(), event.GetUserID())
	case types.InteractionTypeAppCommand:
		if strings.TrimSpace(event.GetUserID()) == "" {
			return types.ErrorResponse("unable to determine the invoking user")
		}
		return handleAppCommand(event)
	default:
		return types.ErrorResponse("unsupported Google Chat event type")
	}
}

func handleAppCommand(event *types.Event) *types.Response {
	commandName, ok := getCommandName(event.AppCommandID())
	if !ok {
		return types.ErrorResponse("unknown Google Chat command id; set GOOGLE_CHAT_COMMAND_IDS to map command ids to names")
	}
	args := appCommandArgs(event, commandName)
	if args == "" {
		return handleCommandText(commandName, event.GetUserID())
	}
	return handleCommandText(commandName+" "+args, event.GetUserID())
}

func appCommandArgs(event *types.Event, commandName string) string {
	raw := strings.TrimSpace(event.MessageText())
	if raw == "" {
		return ""
	}

	trimmed := strings.TrimSpace(strings.TrimPrefix(raw, "/"))
	commandPrefix := strings.TrimSpace(commandName)
	if commandPrefix == "" {
		return trimmed
	}

	lowerTrimmed := strings.ToLower(trimmed)
	lowerCommand := strings.ToLower(commandPrefix)
	if lowerTrimmed == lowerCommand {
		return ""
	}
	if strings.HasPrefix(lowerTrimmed, lowerCommand+" ") {
		return strings.TrimSpace(trimmed[len(commandPrefix):])
	}

	return trimmed
}

func handleCommandText(text, userID string) *types.Response {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return types.TextResponse(welcomeText())
	}
	trimmed = strings.TrimPrefix(trimmed, "/")
	name, args := splitCommand(trimmed)
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return types.TextResponse(welcomeText())
	}
	handler, ok := commands.Get(name)
	if !ok {
		return types.TextResponse(welcomeText())
	}
	return handler(args, userID)
}

func splitCommand(text string) (string, string) {
	idx := strings.IndexByte(text, ' ')
	if idx < 0 {
		return text, ""
	}
	return text[:idx], strings.TrimSpace(text[idx+1:])
}

func welcomeText() string {
	return strings.Join([]string{
		"Meal Headcount Planner is connected.",
		"Commands:",
		"meal view <date>",
		"meal set <date> <meal_type> <YES|NO>",
		"team-meal view <date>",
		"team-meal set <user_id> <date> <meal_type> <YES|NO>",
		"admin-meal view <user_id> <date>",
		"admin-meal set <user_id> <date> <meal_type> <YES|NO>",
		"admin-meal headcount <date>",
		"work-location view <date>",
		"work-location set <date> <OFFICE|WFH>",
		"meal-type view <date>",
		"meal-type add <date> <meal_type>",
		"day-status view <date>",
		"day-status set <date> <GOVERNMENT_HOLIDAY|OFFICE_CLOSED|SPECIAL_EVENT> [note]",
	}, "\n")
}

func getCommandName(commandID int64) (string, bool) {
	commandIDsOnce.Do(func() {
		commandIDMap = map[int64]string{}
		raw := strings.TrimSpace(os.Getenv("GOOGLE_CHAT_COMMAND_IDS"))
		if raw == "" {
			return
		}

		parsed := map[string]string{}
		if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
			return
		}

		for key, value := range parsed {
			var numericID int64
			if _, err := fmt.Sscan(key, &numericID); err == nil {
				commandIDMap[numericID] = value
			}
		}
	})

	commandName, ok := commandIDMap[commandID]
	return commandName, ok
}
