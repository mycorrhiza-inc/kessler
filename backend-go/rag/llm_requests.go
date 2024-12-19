package rag

import (
	"context"
	"fmt"
	"kessler/search"
	"os"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
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

func createSimpleChatCompletionString(messageRequest MultiplexerChatCompletionRequest) (string, error) {
	modelName := messageRequest.ModelName
	chatHistory := messageRequest.ChatHistory
	client, modelid := createOpenaiClientFromString(modelName)

	// Create message slice for OpenAI request
	var messages []openai.ChatCompletionMessage
	for _, history := range chatHistory {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    string(history.Role),
			Content: history.Content,
		})
	}

	openaiRequest := openai.ChatCompletionRequest{
		Model:     modelid,
		MaxTokens: 2000,
		Messages:  messages,
		Stream:    false,
	}

	ctx := context.Background()
	chatResponse, err := client.CreateChatCompletion(ctx, openaiRequest)
	if err != nil {
		return "", fmt.Errorf("failed to create chat completion: %v", err)
	}
	chatText := chatResponse.Choices[0].Message.Content
	if chatText == "" {
		return "", fmt.Errorf("no chat completion text returned")
	}

	return chatText, nil
}

var test_func_document = openai.FunctionDefinition{
	Name: "get_document_info_from_uuid",
	Parameters: jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"uuid": {
				Type:        jsonschema.String,
				Description: "The UUID of the document",
			},
		},
		Required: []string{"uuid"},
	},
}

func createComplexRequest(messageRequest MultiplexerChatCompletionRequest) (ChatMessage, error) {
	ctx := context.Background()
	modelName := messageRequest.ModelName
	chatHistory := messageRequest.ChatHistory
	if !messageRequest.IsSimpleChat {
		AppendInstructionHeaderToChathistory(&chatHistory)
	}
	client, modelID := createOpenaiClientFromString(modelName)

	// Create message slice for OpenAI request
	var dialogue []openai.ChatCompletionMessage
	for _, history := range chatHistory {
		dialogue = append(dialogue, openai.ChatCompletionMessage{
			Role:    string(history.Role),
			Content: history.Content,
		})
	}

	if len(messageRequest.Functions) > 1 {
		return ChatMessage{}, fmt.Errorf("multiple functions not supported, please fix at some point")
	}
	// Logic of code below must be fixed
	var tools []openai.Tool
	for _, fn := range messageRequest.Functions {
		tools = append(tools, openai.Tool{
			Type:     openai.ToolTypeFunction,
			Function: &fn.Schema,
		})
	}

	fmt.Printf("Asking OpenAI '%v' and providing it %d functions...\n",
		dialogue[0].Content, len(messageRequest.Functions))
	fmt.Printf("Calling OpenAI model %v\n", modelID)
	resp, err := client.CreateChatCompletion(ctx,
		openai.ChatCompletionRequest{
			Model:    modelID,
			Messages: dialogue,
			Tools:    tools,
		},
	)
	fmt.Printf("Finished calling OpenAI model %v\n", modelID)

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
