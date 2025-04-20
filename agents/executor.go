package agents

import (
	"github.com/mateors/llmg/callbacks"
	"github.com/mateors/llmg/chains"
	"github.com/mateors/llmg/schema"
)

const _intermediateStepsOutputKey = "intermediateSteps"

// Executor is the chain responsible for running agents.
type Executor struct {
	Agent            Agent
	Memory           schema.Memory
	CallbacksHandler callbacks.Handler
	ErrorHandler     *ParserErrorHandler

	MaxIterations           int
	ReturnIntermediateSteps bool
}

var (
	_ chains.Chain           = &Executor{}
	_ callbacks.HandlerHaver = &Executor{}
)

// NewExecutor creates a new agent executor with an agent and the tools the agent can use.
func NewExecutor(agent Agent, opts ...Option) *Executor {
	options := executorDefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	return &Executor{
		Agent:                   agent,
		Memory:                  options.memory,
		MaxIterations:           options.maxIterations,
		ReturnIntermediateSteps: options.returnIntermediateSteps,
		CallbacksHandler:        options.callbacksHandler,
		ErrorHandler:            options.errorHandler,
	}
}
