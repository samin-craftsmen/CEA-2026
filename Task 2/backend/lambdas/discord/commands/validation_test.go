package commands

import (
	"testing"

	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/types"
)

func TestDiscordCommandFieldValidation(t *testing.T) {
	tests := []struct {
		name    string
		handler func() *types.InteractionResponse
		want    string
	}{
		{
			name: "meal view invalid date",
			handler: func() *types.InteractionResponse {
				return HandleMeal(discordCommandData("meal", "view", option("date", "2026-13-40")), "12345")
			},
			want: `invalid date "2026-13-40": expected YYYY-MM-DD`,
		},
		{
			name: "meal set invalid meal type",
			handler: func() *types.InteractionResponse {
				return HandleMeal(discordCommandData("meal", "set",
					option("date", "2026-04-24"),
					option("meal_type", "Lunch!"),
					option("status", "YES"),
				), "12345")
			},
			want: `invalid meal_type "Lunch!": use 1-32 lowercase letters, numbers, hyphens, or underscores`,
		},
		{
			name: "team meal set missing employee",
			handler: func() *types.InteractionResponse {
				return HandleTeamMeal(discordCommandData("team-meal", "set",
					option("date", "2026-04-24"),
					option("meal_type", "lunch"),
					option("status", "YES"),
				), "12345")
			},
			want: "employee is required",
		},
		{
			name: "admin meal set invalid status",
			handler: func() *types.InteractionResponse {
				return HandleAdminMeal(discordCommandData("admin-meal", "set",
					option("employee", "99999"),
					option("date", "2026-04-24"),
					option("meal_type", "lunch"),
					option("status", "MAYBE"),
				), "12345")
			},
			want: `invalid status "MAYBE": must be YES or NO`,
		},
		{
			name: "work location set invalid location",
			handler: func() *types.InteractionResponse {
				return HandleWorkLocation(discordCommandData("work-location", "set",
					option("date", "2026-04-24"),
					option("location", "HOME"),
				), "12345")
			},
			want: `invalid location "HOME": must be OFFICE or WFH`,
		},
		{
			name: "day status special event note required",
			handler: func() *types.InteractionResponse {
				return HandleDayStatus(discordCommandData("day-status", "set",
					option("date", "2026-04-24"),
					option("type", "SPECIAL_EVENT"),
				), "12345")
			},
			want: "note is required for SPECIAL_EVENT",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := test.handler()
			if response == nil || response.Data == nil || len(response.Data.Embeds) == 0 {
				t.Fatalf("expected error response embed")
			}
			if got := response.Data.Embeds[0].Description; got != test.want {
				t.Fatalf("error description = %q, want %q", got, test.want)
			}
		})
	}
}

func discordCommandData(name, subcommand string, options ...types.CommandOption) *types.CommandData {
	return &types.CommandData{
		Name: name,
		Options: []types.CommandOption{{
			Name:    subcommand,
			Options: options,
		}},
	}
}

func option(name string, value any) types.CommandOption {
	return types.CommandOption{Name: name, Value: value}
}
