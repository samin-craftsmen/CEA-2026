package handler

import (
	"strings"

	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/commands"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/types"
)

// RouteInteraction dispatches a Discord interaction to the appropriate handler.
func RouteInteraction(interaction *types.Interaction) *types.InteractionResponse {
	if interaction == nil {
		return types.ErrorResponse("Invalid interaction payload")
	}

	switch interaction.Type {
	case types.InteractionTypePing:
		return types.PongResponse()
	case types.InteractionTypeApplicationCommand:
		return handleCommand(interaction)
	case types.InteractionTypeAutocomplete:
		return handleAutocomplete(interaction)
	default:
		return types.ErrorResponse("Unsupported interaction type")
	}
}

func handleCommand(interaction *types.Interaction) *types.InteractionResponse {
	if interaction.Data == nil || strings.TrimSpace(interaction.Data.Name) == "" {
		return types.ErrorResponse("No command data received")
	}

	handler, ok := commands.Get(interaction.Data.Name)
	if !ok {
		return types.ErrorResponse("Unknown command: `/" + interaction.Data.Name + "`")
	}

	userID := interaction.GetUserID()
	if strings.TrimSpace(userID) == "" {
		return types.ErrorResponse("Unable to determine the invoking user")
	}
	return handler(interaction.Data, userID)
}
