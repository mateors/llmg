package llms

type MessageContent struct {
	Role  ChatMessageType
	Parts []ContentPart
}

// ContentPart is an interface all parts of content have to implement.
type ContentPart interface {
	isPart()
}

// TextContent is content with some text.
type TextContent struct {
	Text string
}

func (tc TextContent) String() string {
	return tc.Text
}

func (TextContent) isPart() {}

// ContentResponse is the response returned by a GenerateContent call.
// It can potentially return multiple content choices.
type ContentResponse struct {
	Choices []*ContentChoice
}

// ContentChoice is one of the response choices returned by GenerateContent
// calls.
type ContentChoice struct {
	// Content is the textual content of a response
	Content string

	// StopReason is the reason the model stopped generating output.
	StopReason string

	// GenerationInfo is arbitrary information the model adds to the response.
	GenerationInfo map[string]any

	// FuncCall is non-nil when the model asks to invoke a function/tool.
	// If a model invokes more than one function/tool, this field will only
	// contain the first one.
	FuncCall *FunctionCall

	// ToolCalls is a list of tool calls the model asks to invoke.
	ToolCalls []ToolCall

	// This field is only used with the deepseek-reasoner model and represents the reasoning contents of the assistant message before the final answer.
	ReasoningContent string
}

// FunctionCall is the name and arguments of a function call.
type FunctionCall struct {
	// The name of the function to call.
	Name string `json:"name"`
	// The arguments to pass to the function, as a JSON string.
	Arguments string `json:"arguments"`
}

// ToolCall is a call to a tool (as requested by the model) that should be executed.
type ToolCall struct {
	// ID is the unique identifier of the tool call.
	ID string `json:"id"`
	// Type is the type of the tool call. Typically, this would be "function".
	Type string `json:"type"`
	// FunctionCall is the function call to be executed.
	FunctionCall *FunctionCall `json:"function,omitempty"`
}

func (ToolCall) isPart() {}
