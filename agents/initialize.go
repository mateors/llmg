package agents

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mateors/llmg/callbacks"
	"github.com/mateors/llmg/chains"
	"github.com/mateors/llmg/llms"
	"github.com/mateors/llmg/schema"
	"github.com/mateors/llmg/tools"
)

const _defaultMaxIterations = 5

// AgentType is a string type representing the type of agent to create.
type AgentType string

const (
	// ZeroShotReactDescription is an AgentType constant that represents
	// the "zeroShotReactDescription" agent type.
	ZeroShotReactDescription AgentType = "zeroShotReactDescription"
	// ConversationalReactDescription is an AgentType constant that represents
	// the "conversationalReactDescription" agent type.
	ConversationalReactDescription AgentType = "conversationalReactDescription"
)

// Deprecated: This may be removed in the future; please use NewExecutor instead.
// Initialize is a function that creates a new executor with the specified LLM
// model, tools, agent type, and options. It returns an Executor or an error
// if there is any issues during the creation process.
func Initialize(
	llm llms.Model,
	tools []tools.Tool,
	agentType AgentType,
	opts ...Option,
) (*Executor, error) {
	var agent Agent
	switch agentType {
	case ZeroShotReactDescription:
		agent = NewOneShotAgent(llm, tools, opts...)
	case ConversationalReactDescription:
		agent = NewConversationalAgent(llm, tools, opts...)
	default:
		return &Executor{}, ErrUnknownAgentType
	}
	return NewExecutor(agent, opts...), nil
}

func (e *Executor) Call(ctx context.Context, inputValues map[string]any, _ ...chains.ChainCallOption) (map[string]any, error) { //nolint:lll
	inputs, err := inputsToString(inputValues)
	if err != nil {
		return nil, err
	}
	nameToTool := getNameToTool(e.Agent.GetTools())

	steps := make([]schema.AgentStep, 0)
	for i := 0; i < e.MaxIterations; i++ {
		var finish map[string]any
		steps, finish, err = e.doIteration(ctx, steps, nameToTool, inputs)
		if finish != nil || err != nil {
			return finish, err
		}
	}

	if e.CallbacksHandler != nil {
		e.CallbacksHandler.HandleAgentFinish(ctx, schema.AgentFinish{
			ReturnValues: map[string]any{"output": ErrNotFinished.Error()},
		})
	}
	return e.getReturn(
		&schema.AgentFinish{ReturnValues: make(map[string]any)},
		steps,
	), ErrNotFinished
}

func (e *Executor) doIteration( // nolint
	ctx context.Context,
	steps []schema.AgentStep,
	nameToTool map[string]tools.Tool,
	inputs map[string]string,
) ([]schema.AgentStep, map[string]any, error) {
	actions, finish, err := e.Agent.Plan(ctx, steps, inputs)
	if errors.Is(err, ErrUnableToParseOutput) && e.ErrorHandler != nil {
		formattedObservation := err.Error()
		if e.ErrorHandler.Formatter != nil {
			formattedObservation = e.ErrorHandler.Formatter(formattedObservation)
		}
		steps = append(steps, schema.AgentStep{
			Observation: formattedObservation,
		})
		return steps, nil, nil
	}
	if err != nil {
		return steps, nil, err
	}

	if len(actions) == 0 && finish == nil {
		return steps, nil, ErrAgentNoReturn
	}

	if finish != nil {
		if e.CallbacksHandler != nil {
			e.CallbacksHandler.HandleAgentFinish(ctx, *finish)
		}
		return steps, e.getReturn(finish, steps), nil
	}

	for _, action := range actions {
		steps, err = e.doAction(ctx, steps, nameToTool, action)
		if err != nil {
			return steps, nil, err
		}
	}

	return steps, nil, nil
}

func (e *Executor) doAction(
	ctx context.Context,
	steps []schema.AgentStep,
	nameToTool map[string]tools.Tool,
	action schema.AgentAction,
) ([]schema.AgentStep, error) {
	if e.CallbacksHandler != nil {
		e.CallbacksHandler.HandleAgentAction(ctx, action)
	}

	tool, ok := nameToTool[strings.ToUpper(action.Tool)]
	if !ok {
		return append(steps, schema.AgentStep{
			Action:      action,
			Observation: fmt.Sprintf("%s is not a valid tool, try another one", action.Tool),
		}), nil
	}

	observation, err := tool.Call(ctx, action.ToolInput)
	if err != nil {
		return nil, err
	}

	return append(steps, schema.AgentStep{
		Action:      action,
		Observation: observation,
	}), nil
}

func (e *Executor) getReturn(finish *schema.AgentFinish, steps []schema.AgentStep) map[string]any {
	if e.ReturnIntermediateSteps {
		finish.ReturnValues[_intermediateStepsOutputKey] = steps
	}

	return finish.ReturnValues
}

// GetInputKeys gets the input keys the agent of the executor expects.
// Often "input".
func (e *Executor) GetInputKeys() []string {
	return e.Agent.GetInputKeys()
}

// GetOutputKeys gets the output keys the agent of the executor returns.
func (e *Executor) GetOutputKeys() []string {
	return e.Agent.GetOutputKeys()
}

func (e *Executor) GetMemory() schema.Memory { //nolint:ireturn
	return e.Memory
}

func (e *Executor) GetCallbackHandler() callbacks.Handler { //nolint:ireturn
	return e.CallbacksHandler
}

func inputsToString(inputValues map[string]any) (map[string]string, error) {
	inputs := make(map[string]string, len(inputValues))
	for key, value := range inputValues {
		valueStr, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrExecutorInputNotString, key)
		}

		inputs[key] = valueStr
	}

	return inputs, nil
}

func getNameToTool(t []tools.Tool) map[string]tools.Tool {
	if len(t) == 0 {
		return nil
	}

	nameToTool := make(map[string]tools.Tool, len(t))
	for _, tool := range t {
		nameToTool[strings.ToUpper(tool.Name())] = tool
	}

	return nameToTool
}
