package rag

import (
	"encoding/json"
	"net/http"
)

type RequestBody struct {
	Model       string        `json:"model"`
	ChatHistory []ChatMessage `json:"chat_history"`
}

func HandleBasicChatRequest(w http.ResponseWriter, r *http.Request) {
	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	llmObject := LLMModel{reqBody.Model}

	chatHistory := reqBody.ChatHistory
	chatResponse, err := llmObject.Chat(chatHistory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"response": chatResponse})
}

func HandleRagChatRequest(w http.ResponseWriter, r *http.Request) {
	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	llmObject := LLMModel{reqBody.Model}

	chatHistory := reqBody.ChatHistory
	chatResponse, err := llmObject.RagChat(chatHistory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"response": chatResponse})
}
