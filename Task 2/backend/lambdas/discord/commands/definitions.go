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
	Name         string               `json:"name"`
	Description  string               `json:"description"`
	Type         int                  `json:"type"`
	Required     bool                 `json:"required,omitempty"`
	Autocomplete bool                 `json:"autocomplete,omitempty"`
	Options      []SlashCommandOption `json:"options,omitempty"`
	Choices      []SlashCommandChoice `json:"choices,omitempty"`
}

type SlashCommandChoice struct {
	Name  string `json:"name"`
	Value string `json:"value"`
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
				{
					Name:        "set",
					Description: "Set your meal participation status for a specific date",
					Type:        1, // SUB_COMMAND
					Options: []SlashCommandOption{
						{
							Name:        "date",
							Description: "Date in YYYY-MM-DD format (e.g. 2026-03-15)",
							Type:        3, // STRING
							Required:    true,
						},
						{
							Name:         "meal_type",
							Description:  "The meal to update (e.g. lunch, snacks, dinner)",
							Type:         3, // STRING
							Required:     true,
							Autocomplete: true,
						},
						{
							Name:        "status",
							Description: "YES to opt in, NO to opt out",
							Type:        3, // STRING
							Required:    true,
							Choices: []SlashCommandChoice{
								{Name: "YES — Opt in", Value: "YES"},
								{Name: "NO — Opt out", Value: "NO"},
							},
						},
					},
				},
			},
		},
		{
			Name:        "team-meal",
			Description: "Manage meal participation for your team members (Team Lead only)",
			Options: []SlashCommandOption{
				{
					Name:        "view",
					Description: "View meal participation for your team on a specific date",
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
				{
					Name:        "set",
					Description: "Set a team member's meal participation status",
					Type:        1, // SUB_COMMAND
					Options: []SlashCommandOption{
						{
							Name:        "employee",
							Description: "The team member to update",
							Type:        6, // USER
							Required:    true,
						},
						{
							Name:        "date",
							Description: "Date in YYYY-MM-DD format (e.g. 2026-03-15)",
							Type:        3, // STRING
							Required:    true,
						},
						{
							Name:         "meal_type",
							Description:  "The meal to update (e.g. lunch, snacks, dinner)",
							Type:         3, // STRING
							Required:     true,
							Autocomplete: true,
						},
						{
							Name:        "status",
							Description: "YES to opt in, NO to opt out",
							Type:        3, // STRING
							Required:    true,
							Choices: []SlashCommandChoice{
								{Name: "YES — Opt in", Value: "YES"},
								{Name: "NO — Opt out", Value: "NO"},
							},
						},
					},
				},
			},
		},
		{
			Name:        "admin-meal",
			Description: "Manage meal participation for any employee (Admin only)",
			Options: []SlashCommandOption{
				{
					Name:        "view",
					Description: "View an employee's meal participation for a specific date",
					Type:        1, // SUB_COMMAND
					Options: []SlashCommandOption{
						{
							Name:        "employee",
							Description: "The employee to view",
							Type:        6, // USER
							Required:    true,
						},
						{
							Name:        "date",
							Description: "Date in YYYY-MM-DD format (e.g. 2026-03-15)",
							Type:        3, // STRING
							Required:    true,
						},
					},
				},
				{
					Name:        "set",
					Description: "Set an employee's meal participation status",
					Type:        1, // SUB_COMMAND
					Options: []SlashCommandOption{
						{
							Name:        "employee",
							Description: "The employee to update",
							Type:        6, // USER
							Required:    true,
						},
						{
							Name:        "date",
							Description: "Date in YYYY-MM-DD format (e.g. 2026-03-15)",
							Type:        3, // STRING
							Required:    true,
						},
						{
							Name:         "meal_type",
							Description:  "The meal to update (e.g. lunch, snacks, dinner)",
							Type:         3, // STRING
							Required:     true,
							Autocomplete: true,
						},
						{
							Name:        "status",
							Description: "YES to opt in, NO to opt out",
							Type:        3, // STRING
							Required:    true,
							Choices: []SlashCommandChoice{
								{Name: "YES — Opt in", Value: "YES"},
								{Name: "NO — Opt out", Value: "NO"},
							},
						},
					},
				},
				{
					Name:        "headcount",
					Description: "View headcount summary for a specific date (how many opted in/out per meal type)",
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
		{
			Name:        "work-location",
			Description: "Manage your work location (Office/WFH)",
			Options: []SlashCommandOption{
				{
					Name:        "view",
					Description: "View your work location for a specific date",
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
				{
					Name:        "set",
					Description: "Set your work location for a specific date (before 9pm cutoff)",
					Type:        1, // SUB_COMMAND
					Options: []SlashCommandOption{
						{
							Name:        "date",
							Description: "Date in YYYY-MM-DD format (e.g. 2026-03-15)",
							Type:        3, // STRING
							Required:    true,
						},
						{
							Name:        "location",
							Description: "Your work location",
							Type:        3, // STRING
							Required:    true,
							Choices: []SlashCommandChoice{
								{Name: "🏢 Office", Value: "OFFICE"},
								{Name: "🏠 Work from Home", Value: "WFH"},
							},
						},
					},
				},
			},
		},
		{
			Name:        "meal-type",
			Description: "Manage meal types for a specific date (Admin only)",
			Options: []SlashCommandOption{
				{
					Name:        "view",
					Description: "View available meal types for a specific date",
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
				{
					Name:        "add",
					Description: "Add a new meal type for a specific date",
					Type:        1, // SUB_COMMAND
					Options: []SlashCommandOption{
						{
							Name:        "date",
							Description: "Date in YYYY-MM-DD format (e.g. 2026-03-15)",
							Type:        3, // STRING
							Required:    true,
						},
						{
							Name:        "meal_type",
							Description: "Name of the meal type to add (e.g. dinner, iftar)",
							Type:        3, // STRING
							Required:    true,
						},
					},
				},
			},
		},
		{
			Name:        "day-status",
			Description: "Manage the status of a specific day (Admin only)",
			Options: []SlashCommandOption{
				{
					Name:        "set",
					Description: "Set the status for a specific day",
					Type:        1, // SUB_COMMAND
					Options: []SlashCommandOption{
						{
							Name:        "date",
							Description: "Date in YYYY-MM-DD format (e.g. 2026-03-15)",
							Type:        3, // STRING
							Required:    true,
						},
						{
							Name:        "type",
							Description: "Type of day status",
							Type:        3, // STRING
							Required:    true,
							Choices: []SlashCommandChoice{
								{Name: "🏛️ Government Holiday", Value: "GOVERNMENT_HOLIDAY"},
								{Name: "🔒 Office Closed", Value: "OFFICE_CLOSED"},
								{Name: "🎉 Special Event", Value: "SPECIAL_EVENT"},
							},
						},
						{
							Name:        "note",
							Description: "Additional note (required for Special Event)",
							Type:        3, // STRING
							Required:    false,
						},
					},
				},
				{
					Name:        "view",
					Description: "View the status of a specific day",
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
