package rag

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

var openaiKey = os.Getenv("OPENAI_API_KEY")

func createOpenaiClientFromString(model_name string) (*openai.Client, string) {
	switch model_name {
	case "gpt-4o", "gpt-4o-mini":
		return openai.NewClient(openaiKey), openai.GPT4oLatest
	default:
		// Return openai for now, refactor later to deal with stuff
		return openai.NewClient(openaiKey), openai.GPT4oLatest
		// panic(fmt.Sprintf("Unsupported model name: %s", model_name))
	}
}

type FunctionCall struct {
	schema   openai.FunctionDefinition
	function func(string) (ToolCallResults, error)
}
type ToolCallResults struct {
	response  string
	citations *[]SearchData
}

type MultiplexerChatCompletionRequest struct {
	modelName   string
	chatHistory []ChatMessage
	functions   []FunctionCall
}

func createSimpleChatCompletionString(messageRequest MultiplexerChatCompletionRequest) (string, error) {
	modelName := messageRequest.modelName
	chatHistory := messageRequest.chatHistory
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
	modelName := messageRequest.modelName
	chatHistory := messageRequest.chatHistory
	client, modelID := createOpenaiClientFromString(modelName)

	// Create message slice for OpenAI request
	var dialogue []openai.ChatCompletionMessage
	for _, history := range chatHistory {
		dialogue = append(dialogue, openai.ChatCompletionMessage{
			Role:    string(history.Role),
			Content: history.Content,
		})
	}

	if len(messageRequest.functions) > 1 {
		return ChatMessage{}, fmt.Errorf("multiple functions not supported, please fix at some point")
	}
	// Logic of code below must be fixed
	var tools []openai.Tool
	for _, fn := range messageRequest.functions {
		tools = append(tools, openai.Tool{
			Type:     openai.ToolTypeFunction,
			Function: &fn.schema,
		})
	}

	fmt.Printf("Asking OpenAI '%v' and providing it %d functions...\n",
		dialogue[0].Content, len(messageRequest.functions))
	resp, err := client.CreateChatCompletion(ctx,
		openai.ChatCompletionRequest{
			Model:    modelID,
			Messages: dialogue,
			Tools:    tools,
		},
	)
	if err != nil || len(resp.Choices) != 1 {
		return ChatMessage{}, fmt.Errorf("completion error: err:%v len(choices):%v", err, len(resp.Choices))
	}
	msg := resp.Choices[0].Message
	if resp.Choices[0].FinishReason != "tool_calls" && msg.Content != "" {
		return ChatMessage{
			msg.Content,
			"assistant",
			nil,
			nil,
		}, nil
	}
	if len(msg.ToolCalls) == 0 {
		return ChatMessage{}, fmt.Errorf("no tool calls in openai response, dispite stopping for a tool completion")
	}
	if len(msg.ToolCalls) > 1 {
		return ChatMessage{}, fmt.Errorf("we only support one function call per message, got %v, please fix", len(msg.ToolCalls))
	}
	// simulate calling the function & responding to OpenAI
	dialogue = append(dialogue, msg)
	fmt.Printf("OpenAI called us back wanting to invoke our function '%v' with params '%v'\n",
		msg.ToolCalls[0].Function.Name, msg.ToolCalls[0].Function.Arguments)
	run_func := messageRequest.functions[0].function
	run_results, err := run_func(msg.ToolCalls[0].Function.Arguments)
	if err != nil {
		return ChatMessage{}, fmt.Errorf("error running function: %v", err)
	}
	dialogue = append(dialogue, openai.ChatCompletionMessage{
		Role:       openai.ChatMessageRoleTool,
		Content:    run_results.response,
		Name:       msg.ToolCalls[0].Function.Name,
		ToolCallID: msg.ToolCalls[0].ID,
	})
	resp, err = client.CreateChatCompletion(ctx,
		openai.ChatCompletionRequest{
			Model:    modelID,
			Messages: dialogue,
			Tools:    tools,
		},
	)
	if err != nil || len(resp.Choices) != 1 {
		return ChatMessage{}, fmt.Errorf("2nd completion error: err:%v len(choices):%v\n", err, len(resp.Choices))
	}

	// display OpenAI's response to the original question utilizing our function
	msg = resp.Choices[0].Message
	if msg.Content == "" {
		return ChatMessage{}, fmt.Errorf("OpenAI returned an empty response")
	}
	return ChatMessage{
		msg.Content,
		"assistant",
		run_results.citations,
		&[]SimpleChatMessage{},
	}, nil
}
