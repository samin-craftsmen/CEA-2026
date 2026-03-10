package commands

// SlashCommandDefinition describes a command to register with the Discord API.
// Register via PUT https://discord.com/api/v10/applications/{APP_ID}/commands
// with Authorization: Bot {TOKEN} header.
type SlashCommandDefinition struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Options     []SlashCommandOption `json:"options,omitempty"`
}

type SlashCommandOption struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Type        int                  `json:"type"`
	Required    bool                 `json:"required,omitempty"`
	Options     []SlashCommandOption `json:"options,omitempty"`
}

// Definitions returns all slash commands that should be registered with Discord.
// Add new command definitions here when adding new command handlers.
func Definitions() []SlashCommandDefinition {
	return []SlashCommandDefinition{
		{
			Name:        "meal",
			Description: "Manage meal participation",
			Options: []SlashCommandOption{
				{
					Name:        "view",
					Description: "View your meal participation for a specific date",
					Type:        1, // SUB_COMMAND
					Options: []SlashCommandOption{
						{
							Name:        "date",
							Description: "Date in YYYY-MM-DD format (e.g. 2026-03-15)",
							Type:        3, // STRING
							Required:    true,
						},
					},
				},
			},
		},
	}
}
