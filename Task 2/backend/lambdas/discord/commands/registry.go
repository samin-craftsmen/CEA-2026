package commands

import "github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/types"

// CommandHandler is the function signature every slash command must implement.
type CommandHandler func(data *types.CommandData, userID string) *types.InteractionResponse

var registry = map[string]CommandHandler{}

// Register adds a command handler to the registry.
// Call this from an init() function in each command file.
func Register(name string, handler CommandHandler) {
	registry[name] = handler
}

// Get retrieves the handler for the given command name.
func Get(name string) (CommandHandler, bool) {
	h, ok := registry[name]
	return h, ok
}
