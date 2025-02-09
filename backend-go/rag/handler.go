package rag

import (
	"encoding/json"
	"fmt"
	"kessler/objects/networking"
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
)

type ChatRequestBody struct {
	Model       string        `json:"model"`
	ChatHistory []ChatMessage `json:"chat_history"`
}

func checkChatAuthorization(token string) (bool, error) {
	if !strings.HasPrefix(token, "Authenticated") {
		return false, nil
	}
	viewerID := strings.TrimPrefix(token, "Authenticated ")
	if viewerID != "anonomous" {
		return true, nil
	}
	return false, nil
}

func HandleBasicChatRequest(w http.ResponseWriter, r *http.Request) {
	var reqBody ChatRequestBody
	// EVERYONE CAN USE BASIC CHAT FOR NOW
	// isAllowed, _ := checkChatAuthorization(r.Header.Get("Authorization"))
	// if !isAllowed {
	// 	http.Error(w, "Cucumber Water For Customer Only", http.StatusPaymentRequired)
	// 	return
	// }
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		errorstring := fmt.Sprintf("Invalid request payload: %v", err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusBadRequest)
		return
	}
	llmObject := LLMModel{reqBody.Model}

	chatHistory := reqBody.ChatHistory
	chatResponse, err := llmObject.Chat(chatHistory)
	if err != nil {
		log.Info("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"message": chatResponse})
}

type AdvancedRagRequestBody struct {
	Model       string                  `json:"model"`
	ChatHistory []ChatMessage           `json:"chat_history"`
	Filters     networking.FilterFields `json:"filters,omitempty"`
}

func HandleRagChatRequest(w http.ResponseWriter, r *http.Request) {
	// EVERYONE CAN USE RAG FOR NOW!
	// isAllowed, _ := checkChatAuthorization(r.Header.Get("Authorization"))
	// if !isAllowed {
	// 	http.Error(w, "Cucumber Water For Customer Only", http.StatusPaymentRequired)
	// 	return
	// }
	var reqBody AdvancedRagRequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		errorstring := fmt.Sprintf("Invalid request payload: %v", err)
		log.Info(errorstring)
		http.Error(w, errorstring, http.StatusBadRequest)
		return
	}
	llmObject := LLMModel{reqBody.Model}

	chatHistory := reqBody.ChatHistory
	filters := reqBody.Filters
	chatResponse, err := llmObject.RagChat(chatHistory, filters)
	if err != nil {
		log.Info("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"message": chatResponse})
}
