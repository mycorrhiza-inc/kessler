package rag

import (
	"context"
	"fmt"
	"kessler/search"
	"os"

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
		fmt.Println("Model not found, using default model: gpt-4o")
		return openai.NewClient(openaiKey), "gpt-4o"
	}
}

type FunctionCall struct {
	Schema openai.FunctionDefinition
	Func   func(string) (ToolCallResults, error)
}
type ToolCallResults struct {
	Response  string
	Citations *[]search.SearchData
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

// TODO: I tried once to simplify and break down this function into component pieces and failed, sorry - nic
func LLMComplexRequest(messageRequest MultiplexerChatCompletionRequest) (ChatMessage, error) {
	if len(messageRequest.Functions) > 1 {
		return ChatMessage{}, fmt.Errorf("multiple functions not supported, please fix at some point")
	}
	ctx := context.Background()
	modelName := messageRequest.ModelName
	chatHistory := messageRequest.ChatHistory
	if !messageRequest.IsSimpleChat {
		AppendInstructionHeaderToChathistory(&chatHistory)
	}
	client, modelID := createOpenaiClientFromString(modelName)
	dialogue := ComplexToOAIMessages(chatHistory)
	tools := FunctionCallsToOAI(messageRequest.Functions)

	fmt.Printf("Calling OpenAI model %v\n", modelID)
	resp, err := client.CreateChatCompletion(ctx,
		openai.ChatCompletionRequest{
			Model:    modelID,
			Messages: dialogue,
			Tools:    tools,
		},
	)
	// fmt.Printf("Finished calling OpenAI model %v\n", modelID)

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
	if len(msg.ToolCalls) > 1 {
		// return ChatMessage{}, fmt.Errorf("we only support one function call per message, got %v, please fix", len(msg.ToolCalls))
		fmt.Printf("Warning: we only support one function call per message, got %v, proceeding to only do the first tool call.\n", len(msg.ToolCalls))
	}
	dialogue = append(dialogue, msg)
	contextMessages := []openai.ChatCompletionMessage{msg}
	returnCitations := []search.SearchData{}
	for toolCallID := range msg.ToolCalls[:1] {
		fmt.Printf("OpenAI called us back wanting to invoke our function '%v' with params '%v'\n",
			msg.ToolCalls[toolCallID].Function.Name, msg.ToolCalls[toolCallID].Function.Arguments)
		run_func := messageRequest.Functions[toolCallID].Func
		run_results, err := run_func(msg.ToolCalls[toolCallID].Function.Arguments)
		if err != nil {
			return ChatMessage{}, fmt.Errorf("error running function: %v", err)
		}
		toolMsg := openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleTool,
			Content:    run_results.Response,
			Name:       msg.ToolCalls[toolCallID].Function.Name,
			ToolCallID: msg.ToolCalls[toolCallID].ID,
		}
		dialogue = append(dialogue, toolMsg)
		contextMessages = append(contextMessages, toolMsg)
		if run_results.Citations != nil {
			returnCitations = append(returnCitations, *run_results.Citations...)
		}

	}
	fmt.Printf("Asked OAI for all tool calls,\n")
	resp, err = client.CreateChatCompletion(ctx,
		openai.ChatCompletionRequest{
			// Removing ability to recursively call tools, it gets one shot for now.
			// Tools:    tools,
			Model:    modelID,
			Messages: dialogue,
			Tools:    []openai.Tool{},
		},
	)
	fmt.Printf("Finished calling OpenAI model %v\n", modelID)
	if err != nil || len(resp.Choices) != 1 {
		fmt.Printf("2nd OpenAI completion error: err:%v len(choices):%v\n", err, len(resp.Choices))
		return ChatMessage{}, fmt.Errorf("2nd completion error: err:%v len(choices):%v", err, len(resp.Choices))
	}

	// display OpenAI's response to the original question utilizing our function
	msg = resp.Choices[0].Message
	if msg.Content == "" {
		fmt.Printf("OpenAI returned an empty response\n")
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
