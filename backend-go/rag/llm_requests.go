package rag

import (
	"context"
	"fmt"
	"kessler/search"
	"os"

	"github.com/charmbracelet/log"
	openai "github.com/sashabaranov/go-openai"
)

var openaiKey = os.Getenv("OPENAI_API_KEY")

func createOpenaiClientFromString(model_name string) (*openai.Client, string) {
	switch model_name {
	case "gpt-4o", "gpt-4o-mini":
		// return openai.NewClient(openaiKey), openai.GPT4oLatest
		// Apparently "gpt-4o" supports function calling but "gpt-4o-latest" doesnt. Good software design fellas
		return openai.NewClient(openaiKey), "gpt-4o"
	default:
		log.Info("Model not found, using default model: gpt-4o")
		return openai.NewClient(openaiKey), "gpt-4o"
	}
}

type FunctionCall struct {
	Schema openai.FunctionDefinition
	Func   func(string) (ToolCallResults, error)
}
type ToolCallResults struct {
	Response  string
	Citations *[]search.SearchDataHydrated
}

type MultiplexerChatCompletionRequest struct {
	ModelName    string
	ChatHistory  []ChatMessage
	Functions    []FunctionCall
	IsSimpleChat bool
}

func FunctionCallsToOAI(funcs []FunctionCall) []openai.Tool {
	var tools []openai.Tool
	for _, fn := range funcs {
		tools = append(tools, openai.Tool{
			Type:     openai.ToolTypeFunction,
			Function: &fn.Schema,
		})
	}
	return tools
}

func generateFunctionDictionary(functionCalls []FunctionCall) map[string]FunctionCall {
	return_map := make(map[string]FunctionCall)
	for _, funcCall := range functionCalls {
		return_map[funcCall.Schema.Name] = funcCall
	}
	return return_map
}

// TODO: I tried once to simplify and break down this function into component pieces and failed, sorry - nic
func LLMComplexRequest(messageRequest MultiplexerChatCompletionRequest) (ChatMessage, error) {
	// max_tool_calls := 2
	// if len(messageRequest.Functions) > 1 {
	// 	return ChatMessage{}, fmt.Errorf("multiple functions not supported, please fix at some point")
	// }
	ctx := context.Background()
	modelName := messageRequest.ModelName
	chatHistory := messageRequest.ChatHistory
	if !messageRequest.IsSimpleChat {
		AppendInstructionHeaderToChathistory(&chatHistory)
	}
	client, modelID := createOpenaiClientFromString(modelName)
	dialogue := ComplexToOAIMessages(chatHistory)
	tools := FunctionCallsToOAI(messageRequest.Functions)
	toolDictionary := generateFunctionDictionary(messageRequest.Functions)

	log.Info(fmt.Sprintf("Calling OpenAI model %v\n", modelID))
	resp, err := client.CreateChatCompletion(ctx,
		openai.ChatCompletionRequest{
			Model:    modelID,
			Messages: dialogue,
			Tools:    tools,
		},
	)
	// log.Info(fmt.Sprintf("Finished calling OpenAI model %v\n", modelID))

	if err != nil || len(resp.Choices) != 1 {
		return ChatMessage{}, fmt.Errorf("completion error: err:%v len(choices):%v", err, len(resp.Choices))
	}
	msg := resp.Choices[0].Message
	if resp.Choices[0].FinishReason != "tool_calls" && msg.Content != "" {
		return ChatMessage{
			Content:   msg.Content,
			Role:      "assistant",
			Citations: nil,
			Context:   nil,
		}, nil
	}
	if len(msg.ToolCalls) == 0 {
		return ChatMessage{}, fmt.Errorf("no tool calls in openai response, dispite stopping for a tool completion")
	}
	// if len(msg.ToolCalls) > 1 {
	// 	// return ChatMessage{}, fmt.Errorf("we only support one function call per message, got %v, please fix", len(msg.ToolCalls))
	// 	log.Info(fmt.Sprintf("Warning: we only support one function call per message, got %v, proceeding to only do the first tool call.\n", len(msg.ToolCalls)))
	// }
	dialogue = append(dialogue, msg)
	contextMessages := []openai.ChatCompletionMessage{msg}
	returnCitations := []search.SearchDataHydrated{}
	for _, toolCall := range msg.ToolCalls {
		// toolCallName := tool_call.Function.Name
		// tool_call_params := tool_call.Function.Arguments
		fmt.Printf("OpenAI called us back wanting to invoke our function '%v' with params '%v'\n",
			toolCall.Function.Name, toolCall.Function.Arguments)
		run_func := toolDictionary[toolCall.Function.Name].Func
		run_results, err := run_func(toolCall.Function.Arguments)
		if err != nil {
			return ChatMessage{}, fmt.Errorf("error running function: %v", err)
		}
		toolMsg := openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleTool,
			Content:    run_results.Response,
			Name:       toolCall.Function.Name,
			ToolCallID: toolCall.ID,
		}
		dialogue = append(dialogue, toolMsg)
		contextMessages = append(contextMessages, toolMsg)
		if run_results.Citations != nil {
			returnCitations = append(returnCitations, *run_results.Citations...)
		}

	}
	log.Info(fmt.Sprintf("Asked OAI for all tool calls,\n"))
	resp, err = client.CreateChatCompletion(ctx,
		openai.ChatCompletionRequest{
			// Removing ability to recursively call tools, it gets one shot for now.
			// Tools:    tools,
			Model:    modelID,
			Messages: dialogue,
			Tools:    []openai.Tool{},
		},
	)
	log.Info(fmt.Sprintf("Finished calling OpenAI model %v\n", modelID))
	if err != nil || len(resp.Choices) != 1 {
		log.Info(fmt.Sprintf("2nd OpenAI completion error: err:%v len(choices):%v\n", err, len(resp.Choices)))
		return ChatMessage{}, fmt.Errorf("2nd completion error: err:%v len(choices):%v", err, len(resp.Choices))
	}

	// display OpenAI's response to the original question utilizing our function
	msg = resp.Choices[0].Message
	if msg.Content == "" {
		log.Info(fmt.Sprintf("OpenAI returned an empty response\n"))
		return ChatMessage{}, fmt.Errorf("OpenAI returned an empty response")
	}
	simple_returns := OAIMessagesToSimples(contextMessages)
	return ChatMessage{
		Content:   msg.Content,
		Role:      "assistant",
		Citations: &returnCitations,
		Context:   &simple_returns,
	}, nil
}
