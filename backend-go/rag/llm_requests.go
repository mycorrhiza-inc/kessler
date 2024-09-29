
import (
	"context"
	"fmt"
	"io"

	openai "github.com/sashabaranov/go-openai"
)

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
		Stream:    true,
	}

	ctx := context.Background()
	stream, err := client.CreateChatCompletionStream(ctx, openaiRequest)
	if err != nil {
		return "", fmt.Errorf("failed to create chat completion stream: %v", err)
	}
	defer stream.Close()

	var chatResponse string
	for {
		response, err := stream.Recv()
		if err != nil {
			if err != io.EOF {
				return "", fmt.Errorf("stream error: %v", err)
			}
			break
		}
		chatResponse += response.Choices[0].Delta.Content
	}

	return chatResponse, nil
}
