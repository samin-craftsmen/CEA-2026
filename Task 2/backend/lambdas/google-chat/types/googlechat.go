package types

// Interaction types from Google Chat.
const (
	InteractionTypeMessage      = "MESSAGE"
	InteractionTypeCardClicked  = "CARD_CLICKED"
	InteractionTypeAddedToSpace = "ADDED_TO_SPACE"
)

// Event is the Google Chat event payload sent to the bot.
type Event struct {
	Type   string  `json:"type"`
	User   User    `json:"user"`
	Space  Space   `json:"space"`
	Action *Action `json:"action,omitempty"`
	// For slash commands the message contains the command text.
	Message *Message `json:"message,omitempty"`
}

type User struct {
	Name        string `json:"name"`        // "users/12345"
	DisplayName string `json:"displayName"`
}

type Space struct {
	Name string `json:"name"`
}

type Message struct {
	Text          string         `json:"text"`
	SlashCommand  *SlashCommand  `json:"slashCommand,omitempty"`
	ArgumentText  string         `json:"argumentText,omitempty"`
}

type SlashCommand struct {
	CommandID int64 `json:"commandId"`
}

// Action is sent when a user clicks a button on a card.
type Action struct {
	ActionMethodName string       `json:"actionMethodName"`
	Parameters       []ActionParam `json:"parameters,omitempty"`
}

type ActionParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Response is the JSON response sent back to Google Chat.
type Response struct {
	Text  string `json:"text,omitempty"`
	Cards []Card `json:"cards,omitempty"`
}

type Card struct {
	Header   *CardHeader   `json:"header,omitempty"`
	Sections []CardSection `json:"sections,omitempty"`
}

type CardHeader struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle,omitempty"`
}

type CardSection struct {
	Widgets []Widget `json:"widgets"`
}

type Widget struct {
	TextParagraph *TextParagraph `json:"textParagraph,omitempty"`
	KeyValue      *KeyValue      `json:"keyValue,omitempty"`
}

type TextParagraph struct {
	Text string `json:"text"`
}

type KeyValue struct {
	TopLabel    string `json:"topLabel,omitempty"`
	Content     string `json:"content"`
	BottomLabel string `json:"bottomLabel,omitempty"`
}

// GetUserID extracts the bare user ID from the "users/12345" name.
func (e *Event) GetUserID() string {
	name := e.User.Name
	const prefix = "users/"
	if len(name) > len(prefix) {
		return name[len(prefix):]
	}
	return name
}

// GetParam returns the value of a named action parameter.
func (e *Event) GetParam(key string) string {
	if e.Action == nil {
		return ""
	}
	for _, p := range e.Action.Parameters {
		if p.Key == key {
			return p.Value
		}
	}
	return ""
}

// --- Response builders ---

func TextResponse(text string) *Response {
	return &Response{Text: text}
}

func CardResponse(title, subtitle string, widgets []Widget) *Response {
	return &Response{
		Cards: []Card{{
			Header:   &CardHeader{Title: title, Subtitle: subtitle},
			Sections: []CardSection{{Widgets: widgets}},
		}},
	}
}

func ErrorResponse(msg string) *Response {
	return TextResponse("❌ Error: " + msg)
}

func KV(label, content string) Widget {
	return Widget{KeyValue: &KeyValue{TopLabel: label, Content: content}}
}

func Para(text string) Widget {
	return Widget{TextParagraph: &TextParagraph{Text: text}}
}
