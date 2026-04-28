package handler

import (
	"sync"
	"testing"

	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/types"
)

func TestRouteEventAppCommandValidatesInvalidDateFromArgumentText(t *testing.T) {
	t.Setenv("GOOGLE_CHAT_COMMAND_IDS", `{"101":"meal"}`)
	resetCommandLookup()

	response := RouteEvent(&types.Event{
		Type: types.InteractionTypeAppCommand,
		Chat: &types.ChatEvent{
			User: &types.User{Name: "users/12345"},
			AppCommandPayload: &types.AppCommandPayload{
				AppCommandMetadata: &types.AppCommandMetadata{AppCommandID: 101},
				Message:            &types.Message{ArgumentText: "view 2026-13-40"},
			},
		},
	})

	if response == nil {
		return
	}
	if got, want := response.Text, `Error: invalid date "2026-13-40": expected YYYY-MM-DD`; got != want {
		t.Fatalf("response.Text = %q, want %q", got, want)
	}
}

func TestRouteEventAppCommandStripsEchoedCommandTextBeforeValidation(t *testing.T) {
	t.Setenv("GOOGLE_CHAT_COMMAND_IDS", `{"101":"meal"}`)
	resetCommandLookup()

	response := RouteEvent(&types.Event{
		Type: types.InteractionTypeAppCommand,
		Chat: &types.ChatEvent{
			User: &types.User{Name: "users/12345"},
			AppCommandPayload: &types.AppCommandPayload{
				AppCommandMetadata: &types.AppCommandMetadata{AppCommandID: 101},
				Message:            &types.Message{Text: "/meal view 2026-13-40"},
			},
		},
	})

	if response == nil {
		return
	}
	if got, want := response.Text, `Error: invalid date "2026-13-40": expected YYYY-MM-DD`; got != want {
		t.Fatalf("response.Text = %q, want %q", got, want)
	}
}

func resetCommandLookup() {
	commandIDsOnce = sync.Once{}
	commandIDMap = nil
}
