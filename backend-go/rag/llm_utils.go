package rag

import (
	"context"
	"kessler/common/objects/networking"
	"kessler/search"

	openai "github.com/sashabaranov/go-openai"
)

// Use custom enums in place of Python's Enum class

// Define the structures in place of Pydantic models

type RAGChat struct {
	Model       string        `json:"model"`
	ChatHistory []ChatMessage `json:"chat_history"`
}

type SimpleChatMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type ChatMessage struct {
	// SimpleChatMessage
	Content   string                       `json:"content"`
	Role      string                       `json:"role"`
	Citations *[]search.SearchDataHydrated `json:"citations,omitempty"`
	Context   *[]SimpleChatMessage         `json:"context,omitempty"`
}

// Function converting simple to ChatMessage
func SimpleToChatMessage(msg SimpleChatMessage) ChatMessage {
	return ChatMessage{
		// SimpleChatMessage: msg,
		Content:   msg.Content,
		Role:      msg.Role,
		Citations: &[]search.SearchDataHydrated{},
		Context:   &[]SimpleChatMessage{},
	}
}

// Function converting ChatMessage to simple openai.ChatCompletionMessage
func ChatMessageToSimple(msg ChatMessage) SimpleChatMessage {
	return SimpleChatMessage{Content: msg.Content, Role: msg.Role}
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

func ComplexToOAIMessages(messages []ChatMessage) []openai.ChatCompletionMessage {
	var oai_messages []openai.ChatCompletionMessage
	for _, message := range messages {
		oai_messages = append(oai_messages, openai.ChatCompletionMessage{
			Role:    string(message.Role),
			Content: message.Content,
		})
	}
	return oai_messages
}

type LLMModel struct {
	ModelName string
}

var DefaultBigLLMModel = LLMModel{ModelName: "gpt-4o"}

func (model_name LLMModel) Chat(chatHistory []ChatMessage) (ChatMessage, error) {
	requestMultiplex := MultiplexerChatCompletionRequest{
		ChatHistory: chatHistory,
		ModelName:   model_name.ModelName,
		Functions:   []FunctionCall{},
	}
	return LLMComplexRequest(requestMultiplex)
}

// TODO: Add this back in when we have a use case for it.
// var rag_func_call_no_filters = rag_func_call_filters(search.Metadata{})

func (model_name LLMModel) RagChat(chatHistory []ChatMessage, filters networking.FilterFields) (ChatMessage, error) {
	ctx := context.Background()
	requestMultiplex := MultiplexerChatCompletionRequest{
		ChatHistory: chatHistory,
		ModelName:   model_name.ModelName,
		Functions: []FunctionCall{
			rag_func_call_filters(filters),
			more_info_func_call(ctx),
		},
	}
	return LLMComplexRequest(requestMultiplex)
}

type LLM interface {
	Chat(chatHistory []ChatMessage) (ChatMessage, error)
	RagChat(chatHistory []ChatMessage, filters networking.FilterFields) (ChatMessage, error)
}

func SimpleInstruct(model LLM, instruction string) (string, error) {
	chat_history := []ChatMessage{
		{
			Role:    "system",
			Content: instruction,
		},
	}
	instruct_message_result, err := model.Chat(chat_history)
	if err != nil {
		return "", err
	}
	instruct_result_string := instruct_message_result.Content
	return instruct_result_string, nil
}

// Wait to add all the llm utils until you understand how to write concurrent code in go more.
