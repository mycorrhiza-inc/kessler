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

func createOpenaiClientFromString(model_name string) func([]ChatHistory) {
	return func(messages []ChatHistory) {
		switch model_name {
			case "gpt-4o"{
				return openai.NewClient(openaiKey) // Replace with your actual token
			}
		}
	}
}

func HandleBasicChatRequest(w http.ResponseWriter, r *http.Request) {
	c := openai.NewClient(openaiKey) // Replace with your actual token
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

	openaiRequest := openai.ChatCompletionRequest{
		Model:     openai.GPT4oLatest,
		MaxTokens: 2000,
		Messages:  messages,
		Stream:    true,
	}

	ctx := context.Background()
	stream, err := c.CreateChatCompletionStream(ctx, openaiRequest)
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
