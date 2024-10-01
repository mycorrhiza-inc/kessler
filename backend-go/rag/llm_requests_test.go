package rag

import (
	"fmt"
	"testing"
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
	result, err := createSimpleChatCompletionString(modelName, chatHistory)
	if err != nil {
		t.Fail()
		fmt.Println("Error:", err)

	}
	fmt.Println("Result:", result)
}
