package llms

import "errors"

type ChatMessageType string

// ErrUnexpectedChatMessageType is returned when a chat message is of an
// unexpected type.
var ErrUnexpectedChatMessageType = errors.New("unexpected chat message type")

const (
	// ChatMessageTypeAI is a message sent by an AI.
	ChatMessageTypeAI ChatMessageType = "ai"
	// ChatMessageTypeHuman is a message sent by a human.
	ChatMessageTypeHuman ChatMessageType = "human"
	// ChatMessageTypeSystem is a message sent by the system.
	ChatMessageTypeSystem ChatMessageType = "system"
	// ChatMessageTypeGeneric is a message sent by a generic user.
	ChatMessageTypeGeneric ChatMessageType = "generic"
	// ChatMessageTypeFunction is a message sent by a function.
	ChatMessageTypeFunction ChatMessageType = "function"
	// ChatMessageTypeTool is a message sent by a tool.
	ChatMessageTypeTool ChatMessageType = "tool"
)

// ChatMessage represents a message in a chat.
type ChatMessage interface {
	// GetType gets the type of the message.
	GetType() ChatMessageType
	// GetContent gets the content of the message.
	GetContent() string
}

// Named is an interface for objects that have a name.
// type Named interface {
// 	GetName() string
// }

// Statically assert that the types implement the interface.
var (
	//_ ChatMessage = AIChatMessage{}
	_ ChatMessage = HumanChatMessage{}
	//_ ChatMessage = SystemChatMessage{}
	//_ ChatMessage = GenericChatMessage{}
	//_ ChatMessage = FunctionChatMessage{}
	//_ ChatMessage = ToolChatMessage{}
)

// HumanChatMessage is a message sent by a human.
type HumanChatMessage struct {
	Content string
}

func (m HumanChatMessage) GetType() ChatMessageType { return ChatMessageTypeHuman }
func (m HumanChatMessage) GetContent() string       { return m.Content }
