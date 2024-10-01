package rag

import (
	openai "github.com/sashabaranov/go-openai"
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
	Model       string        `json:"model"`
	ChatHistory []ChatMessage `json:"chat_history"`
}

type SimpleChatMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type ChatMessage struct {
	SimpleChatMessage
	Citations *[]SearchData        `json:"citations,omitempty"`
	Context   *[]SimpleChatMessage `json:"context,omitempty"`
}

// Function converting simple to ChatMessage
func SimpleToChatMessage(msg SimpleChatMessage) ChatMessage {
	return ChatMessage{
		SimpleChatMessage: msg,
		Citations:         &[]SearchData{},
		Context:           &[]SimpleChatMessage{},
	}
}

// Function converting ChatMessage to simple openai.ChatCompletionMessage
func ChatMessageToSimple(msg ChatMessage) SimpleChatMessage {
	return msg.SimpleChatMessage
}

func AdvancedMessageContent(msg ChatMessage) string {
	return msg.Content
}

func OAIMsgToSimple(oaiMsg openai.ChatCompletionMessage) SimpleChatMessage {
	return SimpleChatMessage{
		Content: oaiMsg.Content,
		Role:    string(oaiMsg.Role),
	}
}

func SimpleToChatMessages(msgs []SimpleChatMessage) []ChatMessage {
	var keMsgs []ChatMessage
	for _, msg := range msgs {
		keMsg := SimpleToChatMessage(msg)
		keMsgs = append(keMsgs, keMsg)
	}
	return keMsgs
}

func ChatMessageToSimples(keMsgs []ChatMessage) []SimpleChatMessage {
	var msgs []SimpleChatMessage
	for _, keMsg := range keMsgs {
		msg := ChatMessageToSimple(keMsg)
		msgs = append(msgs, msg)
	}
	return msgs
}

func OAIMessagesToSimples(oaiMsgs []openai.ChatCompletionMessage) []SimpleChatMessage {
	var msgs []SimpleChatMessage
	for _, oaiMsg := range oaiMsgs {
		msg := OAIMsgToSimple(oaiMsg)
		msgs = append(msgs, msg)
	}
	return msgs
}

func OAIMessagesToComplex(oaiMsgs []openai.ChatCompletionMessage) []ChatMessage {
	return SimpleToChatMessages(OAIMessagesToSimples(oaiMsgs))
}

type LLMModel struct {
	model_name string
}

func (model_name LLMModel) Chat(chatHistory []ChatMessage) (ChatMessage, error) {
	requestMultiplex := MultiplexerChatCompletionRequest{
		chatHistory: chatHistory,
		modelName:   model_name.model_name,
	}
	return createComplexRequest(requestMultiplex)
}

type LLM interface {
	Chat(chatHistory []ChatMessage) (ChatMessage, error)
}

// Wait to add all the llm utils until you understand how to write concurrent code in go more.
