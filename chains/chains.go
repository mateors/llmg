package chains

import (
	"context"
	"fmt"

	"github.com/mateors/llmg/callbacks"
	"github.com/mateors/llmg/schema"
)

// Chain is the interface all chains must implement.
type Chain interface {
	// Call runs the logic of the chain and returns the output. This method should
	// not be called directly. Use rather the Call, Run or Predict functions that
	// handles the memory and other aspects of the chain.
	Call(ctx context.Context, inputs map[string]any, options ...ChainCallOption) (map[string]any, error)
	// GetMemory gets the memory of the chain.
	GetMemory() schema.Memory
	// GetInputKeys returns the input keys the chain expects.
	GetInputKeys() []string
	// GetOutputKeys returns the output keys the chain returns.
	GetOutputKeys() []string
}

// Call is the standard function used for executing chains.
func Call(ctx context.Context, c Chain, inputValues map[string]any, options ...ChainCallOption) (map[string]any, error) { // nolint: lll
	fullValues := make(map[string]any, 0)
	for key, value := range inputValues {
		fullValues[key] = value
	}

	newValues, err := c.GetMemory().LoadMemoryVariables(ctx, inputValues)
	if err != nil {
		return nil, err
	}

	for key, value := range newValues {
		fullValues[key] = value
	}

	callbacksHandler := getChainCallbackHandler(c)
	if callbacksHandler != nil {
		callbacksHandler.HandleChainStart(ctx, inputValues)
	}

	outputValues, err := callChain(ctx, c, fullValues, options...)
	if err != nil {
		if callbacksHandler != nil {
			callbacksHandler.HandleChainError(ctx, err)
		}
		return outputValues, err
	}

	if callbacksHandler != nil {
		callbacksHandler.HandleChainEnd(ctx, outputValues)
	}

	if err = c.GetMemory().SaveContext(ctx, inputValues, outputValues); err != nil {
		return outputValues, err
	}

	return outputValues, nil
}

func callChain(ctx context.Context, c Chain, fullValues map[string]any, options ...ChainCallOption) (map[string]any, error) {

	if err := validateInputs(c, fullValues); err != nil {
		return nil, err
	}

	outputValues, err := c.Call(ctx, fullValues, options...)
	if err != nil {
		return outputValues, err
	}
	if err := validateOutputs(c, outputValues); err != nil {
		return outputValues, err
	}
	return outputValues, nil
}

// Run can be used to execute a chain if the chain only expects one input and
// one string output.
func Run(ctx context.Context, c Chain, input any, options ...ChainCallOption) (string, error) {

	inputKeys := c.GetInputKeys()
	memoryKeys := c.GetMemory().MemoryVariables(ctx)
	neededKeys := make([]string, 0, len(inputKeys))

	// Remove keys gotten from the memory.
	for _, inputKey := range inputKeys {
		isInMemory := false
		for _, memoryKey := range memoryKeys {
			if inputKey == memoryKey {
				isInMemory = true
				continue
			}
		}
		if isInMemory {
			continue
		}
		neededKeys = append(neededKeys, inputKey)
	}
	if len(neededKeys) != 1 {
		return "", ErrMultipleInputsInRun
	}

	outputKeys := c.GetOutputKeys()
	if len(outputKeys) != 1 {
		return "", ErrMultipleOutputsInRun
	}

	inputValues := map[string]any{neededKeys[0]: input}
	outputValues, err := Call(ctx, c, inputValues, options...)
	if err != nil {
		return "", err
	}

	outputValue, ok := outputValues[outputKeys[0]].(string)
	if !ok {
		return "", ErrWrongOutputTypeInRun
	}
	return outputValue, nil
}

func validateInputs(c Chain, inputValues map[string]any) error {
	for _, k := range c.GetInputKeys() {
		if _, ok := inputValues[k]; !ok {
			return fmt.Errorf("%w: %w: %v", ErrInvalidInputValues, ErrMissingInputValues, k)
		}
	}
	return nil
}

func validateOutputs(c Chain, outputValues map[string]any) error {
	for _, k := range c.GetOutputKeys() {
		if _, ok := outputValues[k]; !ok {
			return fmt.Errorf("%w: %v", ErrInvalidOutputValues, k)
		}
	}
	return nil
}

func getChainCallbackHandler(c Chain) callbacks.Handler {
	if handlerHaver, ok := c.(callbacks.HandlerHaver); ok {
		return handlerHaver.GetCallbackHandler()
	}
	return nil
}
