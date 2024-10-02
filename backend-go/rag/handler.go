package rag

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mycorrhiza-inc/kessler/backend-go/search"
)

type ChatRequestBody struct {
	Model       string        `json:"model"`
	ChatHistory []ChatMessage `json:"chat_history"`
}

func HandleBasicChatRequest(w http.ResponseWriter, r *http.Request) {
	var reqBody ChatRequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	llmObject := LLMModel{reqBody.Model}

	chatHistory := reqBody.ChatHistory
	chatResponse, err := llmObject.Chat(chatHistory)
	if err != nil {
		fmt.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"message": chatResponse})
}

type AdvancedRagRequestBody struct {
	Model       string           `json:"model"`
	ChatHistory []ChatMessage    `json:"chat_history"`
	Filters     *search.Metadata `json:"filters,omitempty"`
}

func HandleRagChatRequest(w http.ResponseWriter, r *http.Request) {
	var reqBody AdvancedRagRequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	llmObject := LLMModel{reqBody.Model}

	chatHistory := reqBody.ChatHistory
	filters := *reqBody.Filters
	chatResponse, err := llmObject.RagChat(chatHistory, filters)
	if err != nil {
		fmt.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"message": chatResponse})
}
