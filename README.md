# LLMG
LLMG is a framework for building AI agents, inspired by LangChain.

## Role Based Chat

```go

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mateors/llmg/llms/ollama"
	"github.com/mateors/llmg/llms"
)

func main() {

	llm, err := ollama.New(ollama.WithModel("llama3.2"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "You are a company branding design wizard."),
		llms.TextParts(llms.ChatMessageTypeHuman, "What would be a good company name a company that makes colorful socks?"),
	}

	// completion, err := llm.GenerateContent(ctx, content, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
	// 	fmt.Print(string(chunk))
	// 	return nil
	// }))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// _ = completion

	cr, err := llm.GenerateContent(ctx, content)
	if err != nil {
		log.Fatal(err)
	}
	for i, cc := range cr.Choices {
		fmt.Println(i, cc)
	}
}

```

## Chat Completion

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mateors/llmg/llms"
	"github.com/mateors/llmg/llms/ollama"
)

func main() {
	llm, err := ollama.New(ollama.WithModel("llama3.2:latest"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	completion, err := llms.GenerateFromSinglePrompt(
		ctx,
		llm,
		"Human: Who was the first man to walk on the moon?\nAssistant:",
		llms.WithTemperature(0.8),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	_ = completion
}

```