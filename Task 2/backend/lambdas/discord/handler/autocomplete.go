package handler

import (
	"strings"

	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/client"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/types"
)

type mealTypesViewReq struct {
	Date string `json:"date"`
}

type mealTypesViewResp struct {
	MealTypes []string `json:"meal_types"`
}

// handleAutocomplete handles Discord autocomplete interactions (type 4).
// Currently supports the meal_type option across all commands by fetching
// the effective meal types for the provided date.
func handleAutocomplete(interaction *types.Interaction) *types.InteractionResponse {
	if interaction.Data == nil || len(interaction.Data.Options) == 0 {
		return types.AutocompleteResponse(defaultChoices())
	}

	// The first option is always the subcommand.
	sub := interaction.Data.Options[0]

	// Extract the date value from sibling options (may be empty if not yet typed).
	var date string
	for _, opt := range sub.Options {
		if opt.Name == "date" {
			if v, ok := opt.Value.(string); ok {
				date = strings.TrimSpace(v)
			}
		}
	}

	if date == "" {
		return types.AutocompleteResponse(defaultChoices())
	}

	var resp mealTypesViewResp
	if err := client.Post("/meal/types/view", mealTypesViewReq{Date: date}, &resp); err != nil || len(resp.MealTypes) == 0 {
		return types.AutocompleteResponse(defaultChoices())
	}

	choices := make([]types.AutocompleteChoice, len(resp.MealTypes))
	for i, mt := range resp.MealTypes {
		choices[i] = types.AutocompleteChoice{
			Name:  strings.ToUpper(mt[:1]) + mt[1:],
			Value: mt,
		}
	}
	return types.AutocompleteResponse(choices)
}

func defaultChoices() []types.AutocompleteChoice {
	return []types.AutocompleteChoice{
		{Name: "Lunch", Value: "lunch"},
		{Name: "Snacks", Value: "snacks"},
	}
}
