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
					Description: "View your meal participation for tomorrow",
					Type:        1, // SUB_COMMAND
				},
			},
		},
	}
}
