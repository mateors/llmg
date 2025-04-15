# LLMG
LLMG is a framework for building AI agents, inspired by LangChain.

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

