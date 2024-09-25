package rag

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

var openaiKey = os.Getenv("OPENAI_API_KEY")

// Define the structure of our request JSON
type ChatHistory struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Model       string        `json:"model"`
	ChatHistory []ChatHistory `json:"chat_history"`
}

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

func HandleBasicChatRequest(w http.ResponseWriter, r *http.Request) {
	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Create message slice for OpenAI request
	var messages []openai.ChatCompletionMessage
	for _, history := range reqBody.ChatHistory {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    history.Role,
			Content: history.Content,
		})
	}

	client, modelid := createOpenaiClientFromString("gpt-4o")
	openaiRequest := openai.ChatCompletionRequest{
		Model:     modelid,
		MaxTokens: 2000,
		Messages:  messages,
		Stream:    true,
	}

	ctx := context.Background()
	stream, err := client.CreateChatCompletionStream(ctx, openaiRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create chat completion stream: %v", err), http.StatusInternalServerError)
		return
	}
	defer stream.Close()

	var chatResponse string
	for {
		response, err := stream.Recv()
		if err != nil {
			if err != io.EOF {
				http.Error(w, fmt.Sprintf("Stream error: %v", err), http.StatusInternalServerError)
			}
			break
		}
		chatResponse += response.Choices[0].Delta.Content
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"response": chatResponse})
}
