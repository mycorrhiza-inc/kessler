package rag

import (
	"fmt"
	"kessler/objects/networking"
	"testing"
)

func TestRag(t *testing.T) {
	history := []SimpleChatMessage{
		{
			Content: "Hello, how can I assist you today?",
			Role:    "assistant",
		},
		{
			Content: "Could you please tell me what xcel energy has to do with the marshall fire by looking at the document database?",
			Role:    "user",
		},
	}
	chatHistory := SimpleToChatMessages(history)
	llmObject := LLMModel{ModelName: "gpt-4o"}
	result, err := llmObject.RagChat(chatHistory, networking.FilterFields{})
	if err != nil {
		t.Fatal(err)
	}
	log.Info("Result:", result)
}
