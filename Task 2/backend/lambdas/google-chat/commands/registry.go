package commands

import "github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/types"

type CommandHandler func(args string, userID string) *types.Response

var registry = map[string]CommandHandler{}

func Register(name string, handler CommandHandler) {
	registry[name] = handler
}

func Get(name string) (CommandHandler, bool) {
	h, ok := registry[name]
	return h, ok
}
