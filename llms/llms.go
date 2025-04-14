package llms

import (
	"context"
	"errors"
)

//type LLM = Model

type Model interface {
	GenerateContent(ctx context.Context, messages []MessageContent, options ...CallOption) (*ContentResponse, error)
	//Call(ctx context.Context, prompt string, options ...CallOption) (string, error)
}

// GenerateFromSinglePrompt is a convenience function for calling an LLM with
// a single string prompt, expecting a single string response. It's useful for
// simple, string-only interactions and provides a slightly more ergonomic API
// than the more general [llms.Model.GenerateContent].
func GenerateFromSinglePrompt(ctx context.Context, llm Model, prompt string, options ...CallOption) (string, error) {
	msg := MessageContent{
		Role:  ChatMessageTypeHuman,
		Parts: []ContentPart{TextContent{Text: prompt}},
	}

	resp, err := llm.GenerateContent(ctx, []MessageContent{msg}, options...)
	if err != nil {
		return "", err
	}

	choices := resp.Choices
	if len(choices) < 1 {
		return "", errors.New("empty response from model")
	}
	c1 := choices[0]
	return c1.Content, nil
}
