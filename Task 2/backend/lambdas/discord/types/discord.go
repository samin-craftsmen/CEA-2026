package types

// Interaction types sent by Discord.
const (
	InteractionTypePing               = 1
	InteractionTypeApplicationCommand = 2
)

// Interaction callback types for responses.
const (
	CallbackTypePong                     = 1
	CallbackTypeChannelMessageWithSource = 4
	CallbackTypeDeferredChannelMessage   = 5
)

// Embed colors.
const (
	ColorSuccess = 0x2ECC71
	ColorError   = 0xE74C3C
	ColorInfo    = 0x3498DB
)

// Interaction represents a Discord interaction payload.
type Interaction struct {
	ID        string       `json:"id"`
	Type      int          `json:"type"`
	Data      *CommandData `json:"data,omitempty"`
	GuildID   string       `json:"guild_id,omitempty"`
	ChannelID string       `json:"channel_id,omitempty"`
	Member    *Member      `json:"member,omitempty"`
	User      *User        `json:"user,omitempty"`
	Token     string       `json:"token"`
}

type CommandData struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Options []CommandOption `json:"options,omitempty"`
}

type CommandOption struct {
	Name    string          `json:"name"`
	Type    int             `json:"type"`
	Value   interface{}     `json:"value,omitempty"`
	Options []CommandOption `json:"options,omitempty"`
}

type Member struct {
	User *User  `json:"user"`
	Nick string `json:"nick,omitempty"`
}

type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
}

type InteractionResponse struct {
	Type int                      `json:"type"`
	Data *InteractionResponseData `json:"data,omitempty"`
}

type InteractionResponseData struct {
	Content string  `json:"content,omitempty"`
	Embeds  []Embed `json:"embeds,omitempty"`
	Flags   int     `json:"flags,omitempty"`
}

type Embed struct {
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	Color       int          `json:"color,omitempty"`
	Fields      []EmbedField `json:"fields,omitempty"`
	Footer      *EmbedFooter `json:"footer,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

type EmbedFooter struct {
	Text string `json:"text"`
}

// GetUserID extracts the user ID from the interaction.
// In guilds it's nested under Member; in DMs it's a top-level User.
func (i *Interaction) GetUserID() string {
	if i.Member != nil && i.Member.User != nil {
		return i.Member.User.ID
	}
	if i.User != nil {
		return i.User.ID
	}
	return ""
}

// --- Response builders ---

func PongResponse() *InteractionResponse {
	return &InteractionResponse{Type: CallbackTypePong}
}

func TextResponse(content string) *InteractionResponse {
	return &InteractionResponse{
		Type: CallbackTypeChannelMessageWithSource,
		Data: &InteractionResponseData{Content: content},
	}
}

func EmbedResponse(embeds ...Embed) *InteractionResponse {
	return &InteractionResponse{
		Type: CallbackTypeChannelMessageWithSource,
		Data: &InteractionResponseData{Embeds: embeds},
	}
}

func ErrorResponse(message string) *InteractionResponse {
	return &InteractionResponse{
		Type: CallbackTypeChannelMessageWithSource,
		Data: &InteractionResponseData{
			Embeds: []Embed{{
				Title:       "Error",
				Description: message,
				Color:       ColorError,
			}},
			Flags: 64, // Ephemeral – only visible to the invoking user
		},
	}
}
