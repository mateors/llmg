# Comprehensive Documentation for LLMG (Large Language Model Gateway)

## Overview

The LLMG is a Go package that provides a unified interface for interacting with various Large Language Models (LLMs). It abstracts away the complexities of different LLM providers, allowing developers to use a consistent API regardless of the underlying model service.

## Core Components

The codebase is organized into several key components:

### 1. LLMs Package

**Purpose:** Defines the core interfaces and types for interacting with language models.

**Key Files:**
- `llms.go`: Contains the main Model interface
- `chat_messages.go`: Defines message types and structures
- `generatecontent.go`: Handles content generation responses
- `options.go`: Provides configuration options for LLM calls

The LLMs package serves as the foundation of the LLMG system. It defines a standard interface (`Model`) that all LLM implementations must follow, allowing for consistent interaction regardless of the specific provider.

**Example - Basic Model Interface:**
```go
// Model interface defines the methods that all LLM implementations must provide
type Model interface {
    GenerateContent(ctx context.Context, messages []MessageContent, options ...CallOption) (*ContentResponse, error)
}
```

### 2. Provider Implementations

**Purpose:** Implements the Model interface for specific LLM providers.

**Key Files:**
- `ollama/ollama.go`: Implementation for Ollama models

Each provider implementation handles the specific requirements of its respective API, including authentication, request formatting, and response parsing, while conforming to the standard Model interface.

**Example - Ollama Implementation:**
```go
// LLM is an Ollama LLM implementation
type LLM struct {
    CallbacksHandler callbacks.Handler
    client           *ollamaclient.Client
    options          options
}

// GenerateContent implements the Model interface
func (o *LLM) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (*llms.ContentResponse, error) {
    // Implementation details...
}
```

### 3. Callbacks System

**Purpose:** Provides hooks for monitoring and responding to events during LLM operations.

**Key Files:**
- `callbacks/callbacks.go`: Defines the Handler interface
- `callbacks/manager.go`: Implements callback management

The callbacks system allows applications to hook into various events in the LLM lifecycle, such as when a request starts, when a response is received, or when an error occurs.

**Example - Callback Handler Interface:**
```go
// Handler is the interface that allows for hooking into specific parts of an LLM application
type Handler interface {
    HandleText(ctx context.Context, text string)
    HandleLLMGenerateContentStart(ctx context.Context, ms []llms.MessageContent)
    HandleLLMGenerateContentEnd(ctx context.Context, res *llms.ContentResponse)
    HandleLLMError(ctx context.Context, err error)
    // Additional methods (currently commented out)...
}
```

### 4. Schema Package

**Purpose:** Defines data structures for agent-based operations.

**Key Files:**
- `schema/schema.go`: Contains structures for agent actions and steps

The schema package provides structures for more complex LLM applications, particularly those involving agents that can take actions based on LLM outputs.

**Example - Agent Structures:**
```go
// AgentAction is the agent's action to take
type AgentAction struct {
    Tool      string
    ToolInput string
    Log       string
    ToolID    string
}

// AgentFinish is the agent's return value
type AgentFinish struct {
    ReturnValues map[string]any
    Log          string
}
```

## Key Workflows

### 1. Basic Text Generation

The most common workflow is generating text from a prompt:

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/mateors/llmg/llms"
    "github.com/mateors/llmg/llms/ollama"
)

func main() {
    // Initialize the Ollama model
    model, err := ollama.New(
        ollama.WithModel("llama2"),
        ollama.WithServerURL("http://localhost:11434"),
    )
    if err != nil {
        panic(err)
    }
    
    // Create a message
    messages := []llms.MessageContent{
        llms.TextParts(llms.ChatMessageTypeHuman, "Explain quantum computing in simple terms"),
    }
    
    // Generate content
    response, err := model.GenerateContent(
        context.Background(),
        messages,
        llms.WithTemperature(0.7),
    )
    if err != nil {
        panic(err)
    }
    
    // Print the response
    fmt.Println(response.Choices[0].Content)
}
```

### 2. Using Callbacks

Callbacks allow you to monitor and respond to events during the LLM operation:

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/mateors/llmg/callbacks"
    "github.com/mateors/llmg/llms"
    "github.com/mateors/llmg/llms/ollama"
)

// Custom callback handler
type LoggingHandler struct{}

func (h *LoggingHandler) HandleText(ctx context.Context, text string) {
    fmt.Println("Text:", text)
}

func (h *LoggingHandler) HandleLLMGenerateContentStart(ctx context.Context, ms []llms.MessageContent) {
    fmt.Println("Generation started with", len(ms), "messages")
}

func (h *LoggingHandler) HandleLLMGenerateContentEnd(ctx context.Context, res *llms.ContentResponse) {
    fmt.Println("Generation completed with", len(res.Choices), "choices")
}

func (h *LoggingHandler) HandleLLMError(ctx context.Context, err error) {
    fmt.Println("Error occurred:", err)
}

func main() {
    // Initialize the model
    model, err := ollama.New(
        ollama.WithModel("llama2"),
        ollama.WithServerURL("http://localhost:11434"),
    )
    if err != nil {
        panic(err)
    }
    
    // Set the callback handler
    model.CallbacksHandler = &LoggingHandler{}
    
    // Create a message and generate content
    // ... (same as previous example)
}
```

### 3. Streaming Responses

For applications that need real-time responses:

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/mateors/llmg/llms"
    "github.com/mateors/llmg/llms/ollama"
)

func main() {
    // Initialize the model
    model, err := ollama.New(
        ollama.WithModel("llama2"),
        ollama.WithServerURL("http://localhost:11434"),
    )
    if err != nil {
        panic(err)
    }
    
    // Create a message
    messages := []llms.MessageContent{
        llms.TextParts(llms.ChatMessageTypeHuman, "Write a short story about AI"),
    }
    
    // Define a streaming function
    streamingFunc := func(ctx context.Context, chunk []byte) error {
        fmt.Print(string(chunk))
        return nil
    }
    
    // Generate content with streaming
    _, err = model.GenerateContent(
        context.Background(),
        messages,
        llms.WithStreamingFunc(streamingFunc),
    )
    if err != nil {
        panic(err)
    }
}
```

## Detailed Component Explanations

### Message Structure

Messages in LLMG follow a specific structure:

1. **MessageContent**: Contains a role (who is sending the message) and parts (the content)
2. **ContentPart**: An interface for different types of content (text, binary, etc.)
3. **TextContent**: A simple text content implementation

```go
// Creating a system message
systemMsg := llms.TextParts(llms.ChatMessageTypeSystem, "You are a helpful assistant.")

// Creating a user message
userMsg := llms.TextParts(llms.ChatMessageTypeHuman, "What is quantum computing?")

// Combining messages for a request
messages := []llms.MessageContent{systemMsg, userMsg}
```

### Call Options

LLMG provides a flexible options system for configuring LLM requests:

```go
// Basic options
response, err := model.GenerateContent(
    context.Background(),
    messages,
    llms.WithModel("gpt-3.5-turbo"),
    llms.WithTemperature(0.7),
)

// More advanced options
options := llms.CallOptions{
    Model:       "gpt-4",
    MaxTokens:   1000,
    Temperature: 0.8,
    TopP:        0.95,
    StopWords:   []string{"STOP", "END"},
}

response, err := model.GenerateContent(
    context.Background(),
    messages,
    llms.WithOptions(options),
)
```

### Callback Manager

For applications that need to manage multiple callbacks:

```go
// Create a callback manager
manager := callbacks.NewCallbackManager()

// Add handlers
manager.AddHandler(&LoggingHandler{})
manager.AddHandler(&MetricsHandler{})

// Use the manager with a model
model.CallbacksHandler = manager
```

## Best Practices

1. **Error Handling**: Always check for errors when calling LLM methods
2. **Context Management**: Use context for timeout and cancellation control
3. **Callback Usage**: Implement callbacks for monitoring and debugging
4. **Option Configuration**: Use appropriate options for your specific use case

## Conclusion

The LLMG package provides a flexible and powerful abstraction for working with various LLM providers. By standardizing the interface and handling provider-specific details internally, it allows developers to focus on building applications rather than managing the complexities of different LLM APIs.

The modular design makes it easy to add support for new providers while maintaining a consistent experience for application developers. The callback system provides powerful hooks for monitoring and responding to LLM operations, making it suitable for a wide range of applications.
