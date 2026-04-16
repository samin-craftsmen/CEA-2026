package handler

import (
	"strings"

	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/commands"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/types"
)

// RouteEvent dispatches a Google Chat event to the appropriate command handler.
func RouteEvent(event *types.Event) *types.Response {
	switch event.Type {
	case types.InteractionTypeAddedToSpace:
		return types.TextResponse("👋 Meal Headcount Planner bot added! Use `/meal`, `/team-meal`, `/admin-meal`, `/meal-type`, `/work-location`, or `/day-status`.")
	case types.InteractionTypeMessage:
		return handleMessage(event)
	default:
		return types.ErrorResponse("unsupported interaction type")
	}
}

func handleMessage(event *types.Event) *types.Response {
	if event.Message == nil {
		return types.ErrorResponse("no message received")
	}

	text := strings.TrimSpace(event.Message.Text)
	if !strings.HasPrefix(text, "/") {
		return nil // ignore non-command messages
	}

	// Split "/command args" → command name + rest
	text = text[1:] // strip leading /
	idx := strings.IndexByte(text, ' ')
	var name, args string
	if idx < 0 {
		name = text
	} else {
		name = text[:idx]
		args = strings.TrimSpace(text[idx+1:])
	}

	h, ok := commands.Get(name)
	if !ok {
		return types.ErrorResponse("unknown command `/" + name + "`")
	}

	userID := event.GetUserID()
	return h(args, userID)
}
