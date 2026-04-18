package types

import "strings"

const (
	InteractionTypeMessage      = "MESSAGE"
	InteractionTypeAddedToSpace = "ADDED_TO_SPACE"
	InteractionTypeAppCommand   = "APP_COMMAND"
)

type Event struct {
	Type              string             `json:"type,omitempty"`
	CommonEventObject *CommonEventObject `json:"commonEventObject,omitempty"`
	Chat              *ChatEvent         `json:"chat,omitempty"`
	Message           *Message           `json:"message,omitempty"`
	User              *User              `json:"user,omitempty"`
	Space             *Space             `json:"space,omitempty"`
	Token             string             `json:"token,omitempty"`
}

type CommonEventObject struct {
	HostApp  string `json:"hostApp,omitempty"`
	Platform string `json:"platform,omitempty"`
}

type ChatEvent struct {
	User                *User                `json:"user,omitempty"`
	EventTime           string               `json:"eventTime,omitempty"`
	MessagePayload      *MessagePayload      `json:"messagePayload,omitempty"`
	AddedToSpacePayload *AddedToSpacePayload `json:"addedToSpacePayload,omitempty"`
	AppCommandPayload   *AppCommandPayload   `json:"appCommandPayload,omitempty"`
}

type MessagePayload struct {
	Space   *Space   `json:"space,omitempty"`
	Message *Message `json:"message,omitempty"`
}

type AddedToSpacePayload struct {
	Space *Space `json:"space,omitempty"`
}

type AppCommandPayload struct {
	AppCommandMetadata *AppCommandMetadata `json:"appCommandMetadata,omitempty"`
	Space              *Space              `json:"space,omitempty"`
	Message            *Message            `json:"message,omitempty"`
}

type AppCommandMetadata struct {
	AppCommandID int64 `json:"appCommandId,omitempty"`
}

type Space struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

type Message struct {
	Text          string `json:"text,omitempty"`
	ArgumentText  string `json:"argumentText,omitempty"`
	FormattedText string `json:"formattedText,omitempty"`
}

type User struct {
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

type Response struct {
	Text string `json:"text,omitempty"`
}

func (e *Event) InteractionType() string {
	if strings.TrimSpace(e.Type) != "" {
		return e.Type
	}
	if e.Chat == nil {
		return ""
	}
	switch {
	case e.Chat.MessagePayload != nil:
		return InteractionTypeMessage
	case e.Chat.AddedToSpacePayload != nil:
		return InteractionTypeAddedToSpace
	case e.Chat.AppCommandPayload != nil:
		return InteractionTypeAppCommand
	default:
		return ""
	}
}

func (e *Event) MessageText() string {
	if e.Message != nil {
		if text := firstNonEmpty(e.Message.ArgumentText, e.Message.Text, e.Message.FormattedText); text != "" {
			return text
		}
	}
	if e.Chat == nil {
		return ""
	}
	if e.Chat.MessagePayload != nil && e.Chat.MessagePayload.Message != nil {
		return firstNonEmpty(
			e.Chat.MessagePayload.Message.ArgumentText,
			e.Chat.MessagePayload.Message.Text,
			e.Chat.MessagePayload.Message.FormattedText,
		)
	}
	if e.Chat.AppCommandPayload != nil && e.Chat.AppCommandPayload.Message != nil {
		return firstNonEmpty(
			e.Chat.AppCommandPayload.Message.ArgumentText,
			e.Chat.AppCommandPayload.Message.Text,
			e.Chat.AppCommandPayload.Message.FormattedText,
		)
	}
	return ""
}

func (e *Event) GetUserID() string {
	if e.User != nil {
		if id := normalizeUserID(e.User.Name); id != "" {
			return id
		}
	}
	if e.Chat != nil && e.Chat.User != nil {
		if id := normalizeUserID(e.Chat.User.Name); id != "" {
			return id
		}
	}
	return ""
}

func (e *Event) HostAppIsChat() bool {
	return e.CommonEventObject != nil && strings.EqualFold(e.CommonEventObject.HostApp, "CHAT")
}

func (e *Event) AppCommandID() int64 {
	if e.Chat == nil || e.Chat.AppCommandPayload == nil || e.Chat.AppCommandPayload.AppCommandMetadata == nil {
		return 0
	}
	return e.Chat.AppCommandPayload.AppCommandMetadata.AppCommandID
}

func TextResponse(text string) *Response {
	return &Response{Text: text}
}

func ErrorResponse(message string) *Response {
	return TextResponse("Error: " + message)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func normalizeUserID(value string) string {
	trimmed := strings.TrimSpace(value)
	trimmed = strings.TrimPrefix(trimmed, "users/")
	trimmed = strings.TrimPrefix(trimmed, "<users/")
	trimmed = strings.TrimSuffix(trimmed, ">")
	return strings.TrimSpace(trimmed)
}
