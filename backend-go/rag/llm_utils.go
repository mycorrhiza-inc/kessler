package rag

import (
	"fmt"
)

// Use custom enums in place of Python's Enum class

// Define the structures in place of Pydantic models
type SearchData struct {
	Name     string `json:"name"`
	Text     string `json:"text"`
	DocID    string `json:"doc_id"`
	SourceID string `json:"source_id"`
}

type RAGChat struct {
	Model       string          `json:"model"`
	ChatHistory []KeChatMessage `json:"chat_history"`
}

// write two functions, one that converts a simple chatmessage into a KeChatMessage with emtpy values for context and citations, while validating the chatRole. Then write another that just throws that info away and turns a KeChatMessage into a simple chatmessage.
type ChatRole string

const (
	User      ChatRole = "user"
	System    ChatRole = "system"
	Assistant ChatRole = "assistant"
)

type SimpleChatMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type KeChatMessage struct {
	Content   string               `json:"content"`
	Role      ChatRole             `json:"role"`
	Citations *[]SearchData        `json:"citations,omitempty"`
	Context   *[]SimpleChatMessage `json:"context,omitempty"`
}

func ValidateChatRole(role string) (ChatRole, error) {
	chatRole := ChatRole(role)
	switch chatRole {
	case User, System, Assistant:
		return chatRole, nil
	default:
		return "", fmt.Errorf("invalid role: %s", role)
	}
}

func SimpleToKeChatMessage(msg SimpleChatMessage) (KeChatMessage, error) {
	role, err := ValidateChatRole(msg.Role)
	if err != nil {
		return KeChatMessage{}, err
	}
	return KeChatMessage{
		Content:   msg.Content,
		Role:      role,
		Citations: &[]SearchData{},
		Context:   &[]SimpleChatMessage{},
	}, nil
}

func KeToSimpleChatMessage(keMsg KeChatMessage) SimpleChatMessage {
	return SimpleChatMessage{
		Content: keMsg.Content,
		Role:    string(keMsg.Role),
	}
}

type LLM interface {
	Achat(chatHistory []map[string]interface{}) (string, error)
}

type KeLLMUtils struct {
	Llm LLM
}

// Mock implementation of LLM interface
type MockLLM struct{}

func (m *MockLLM) Achat(chatHistory []map[string]interface{}) (string, error) {
	return "assistant: Hello, world", nil
}

// Making 'achat' synchronous for simplicity
func (k *KeLLMUtils) Achat(chatHistory []KeChatMessage) (KeChatMessage, error) {
	llamaChatHistory := sanitizeChatHistoryLlamaindex(chatHistory)
	response, err := k.Llm.Achat(llamaChatHistory)
	if err != nil {
		return KeChatMessage{}, err
	}

	strResponse := removePrefixes(response)
	return KeChatMessage{
		Role:    Assistant,
		Content: strResponse,
	}, nil
}
