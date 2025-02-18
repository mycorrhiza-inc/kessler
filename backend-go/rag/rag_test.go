package rag

import (
	"context"
	"kessler/common/llm_utils"
	"kessler/common/objects/networking"
	"testing"

	"github.com/charmbracelet/log"
)

func TestRag(t *testing.T) {
	ctx := context.Background()
	history := []llm_utils.SimpleChatMessage{
		{
			Content: "Hello, how can I assist you today?",
			Role:    "assistant",
		},
		{
			Content: "Could you please tell me what xcel energy has to do with the marshall fire by looking at the document database?",
			Role:    "user",
		},
	}
	chatHistory := llm_utils.SimpleToChatMessages(history)
	llmObject := llm_utils.DefaultBigLLMModel
	ragLLM := RagLLMModel(llmObject)
	result, err := ragLLM.RagChat(ctx, chatHistory, networking.FilterFields{})
	if err != nil {
		t.Fatal(err)
	}
	log.Info("Result:", result)
}
