package rag

import (
	"fmt"
	"testing"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

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
		fmt.Println("Error:", err)
	}
	fmt.Println("Result:", result)
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
