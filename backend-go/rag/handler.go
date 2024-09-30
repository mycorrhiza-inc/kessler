package rag

import (
	"encoding/json"
	"net/http"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

var openaiKey = os.Getenv("OPENAI_API_KEY")

type RequestBody struct {
	Model       string          `json:"model"`
	ChatHistory []ChatMessage `json:"chat_history"`
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

	chatResponse, err := CreateKeChatCompletion(reqBody.Model, reqBody.ChatHistory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"response": chatResponse})
}
