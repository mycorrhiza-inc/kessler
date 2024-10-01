package rag

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
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

func createSimpleChatCompletionString(modelName string, chatHistory []SimpleChatMessage) (string, error) {
	client, modelid := createOpenaiClientFromString(modelName)

	// Create message slice for OpenAI request
	var messages []openai.ChatCompletionMessage
	for _, history := range chatHistory {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    history.Role,
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
