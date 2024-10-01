package rag

import (
	"fmt"
	"testing"

	openai "github.com/sashabaranov/go-openai"
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
			Content: "I'm struggling with my math homework on quadratic equations.",
			Role:    "user",
		},
	}
	multiplex_request := MultiplexerChatCompletionRequest{
		modelName,
		chatHistory,
		[]openai.FunctionDefinition{},
	}
	result, err := createSimpleChatCompletionString(multiplex_request)
	if err != nil {
		t.Fail()
		fmt.Println("Error:", err)
	}
	fmt.Println("Result:", result)
}

func TestChatFunctionCalling(t *testing.T) {
	modelName := "gpt-4o"
	chatHistory := []SimpleChatMessage{
		{
			Content: "Hello, how can I assist you today?",
			Role:    "assistant",
		},
		{
			Content: "Could you please tell me the current weather in Denver CO?",
			Role:    "user",
		},
	}
	multiplex_request := MultiplexerChatCompletionRequest{
		modelName,
		chatHistory,
		[]openai.FunctionDefinition{},
	}
	result, err := createSimpleChatCompletionString(multiplex_request)
	if err != nil {
		t.Fail()
		fmt.Println("Error:", err)
	}
	fmt.Println("Result:", result)
}
