package chains

import (
	"context"

	"github.com/mateors/llmg/callbacks"
)

// ChainCallOption is a function that can be used to modify the behavior of the Call function.
type ChainCallOption func(*chainCallOption)

// For issue #626, each field here has a boolean "set" flag so we can
// distinguish between the case where the option was actually set explicitly
// on chainCallOption, or asked to remain default. The reason we need this is
// that in translating options from ChainCallOption to llms.CallOption, the
// notion of "default value the user didn't explicitly ask to change" is
// violated.
// These flags are hopefully a temporary backwards-compatible solution, until
// we find a more fundamental solution for #626.
type chainCallOption struct {
	// Model is the model to use in an LLM call.
	Model    string
	modelSet bool

	// MaxTokens is the maximum number of tokens to generate to use in an LLM call.
	MaxTokens    int
	maxTokensSet bool

	// Temperature is the temperature for sampling to use in an LLM call, between 0 and 1.
	Temperature    float64
	temperatureSet bool

	// StopWords is a list of words to stop on to use in an LLM call.
	StopWords    []string
	stopWordsSet bool

	// StreamingFunc is a function to be called for each chunk of a streaming response.
	// Return an error to stop streaming early.
	StreamingFunc func(ctx context.Context, chunk []byte) error

	// TopK is the number of tokens to consider for top-k sampling in an LLM call.
	TopK    int
	topkSet bool

	// TopP is the cumulative probability for top-p sampling in an LLM call.
	TopP    float64
	toppSet bool

	// Seed is a seed for deterministic sampling in an LLM call.
	Seed    int
	seedSet bool

	// MinLength is the minimum length of the generated text in an LLM call.
	MinLength    int
	minLengthSet bool

	// MaxLength is the maximum length of the generated text in an LLM call.
	MaxLength    int
	maxLengthSet bool

	// RepetitionPenalty is the repetition penalty for sampling in an LLM call.
	RepetitionPenalty    float64
	repetitionPenaltySet bool

	// CallbackHandler is the callback handler for Chain
	CallbackHandler callbacks.Handler
}

// WithModel is an option for LLM.Call.
func WithModel(model string) ChainCallOption {
	return func(o *chainCallOption) {
		o.Model = model
		o.modelSet = true
	}
}
