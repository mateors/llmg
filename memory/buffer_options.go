package memory

// ConversationBufferOption is a function for creating new buffer
// with other than the default values.
type ConversationBufferOption func(b *ConversationBuffer)

func applyBufferOptions(opts ...ConversationBufferOption) *ConversationBuffer {
	m := &ConversationBuffer{
		ReturnMessages: false,
		InputKey:       "",
		OutputKey:      "",
		HumanPrefix:    "Human",
		AIPrefix:       "AI",
		MemoryKey:      "history",
	}

	for _, opt := range opts {
		opt(m)
	}

	if m.ChatHistory == nil {
		m.ChatHistory = NewChatMessageHistory()
	}

	return m
}
