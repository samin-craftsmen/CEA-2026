package handler

import (
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/commands"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/types"
)

// RouteInteraction dispatches a Discord interaction to the appropriate handler.
func RouteInteraction(interaction *types.Interaction) *types.InteractionResponse {
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
	if interaction.Data == nil {
		return types.ErrorResponse("No command data received")
	}

	handler, ok := commands.Get(interaction.Data.Name)
	if !ok {
		return types.ErrorResponse("Unknown command: `/" + interaction.Data.Name + "`")
	}

	userID := interaction.GetUserID()
	return handler(interaction.Data, userID)
}
