package rag

import (
	"context"
	"fmt"
	"testing"

	"github.com/charmbracelet/log"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

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

func TestSimpleChatCompletionString(t *testing.T) {
	modelName := "gpt-4o"
	chatHistory := []SimpleChatMessage{
		{
			Content: "Hello, how can I assist you today?",
			Role:    "assistant",
		},
		{
			Content: "Can you help me with my homework?",
			Role:    "user",
		},
		{
			Content: "Of course! What subject are you working on?",
			Role:    "assistant",
		},
		{
			Content: "I'm struggling with my math homework deriving a proof that the dirichlet function is discontinuous everywhere.",
			Role:    "user",
		},
	}
	multiplex_request := MultiplexerChatCompletionRequest{
		ModelName:    modelName,
		ChatHistory:  SimpleToChatMessages(chatHistory),
		Functions:    []FunctionCall{},
		IsSimpleChat: true,
	}
	result, err := createSimpleChatCompletionString(multiplex_request)
	if err != nil {
		t.Fail()
		log.Info("Error:", err)
	}
	log.Info("Result:", result)
}

var test_document_func_schema = openai.FunctionDefinition{
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
