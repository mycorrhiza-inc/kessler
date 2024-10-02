package rag

import (
	"encoding/json"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"

	"github.com/mycorrhiza-inc/kessler/backend-go/search"
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
	Content   string               `json:"content"`
	Role      string               `json:"role"`
	Citations *[]search.SearchData `json:"citations,omitempty"`
	Context   *[]SimpleChatMessage `json:"context,omitempty"`
}

// Function converting simple to ChatMessage
func SimpleToChatMessage(msg SimpleChatMessage) ChatMessage {
	return ChatMessage{
		// SimpleChatMessage: msg,
		Content:   msg.Content,
		Role:      msg.Role,
		Citations: &[]search.SearchData{},
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

type LLMModel struct {
	ModelName string
}

func (model_name LLMModel) Chat(chatHistory []ChatMessage) (ChatMessage, error) {
	requestMultiplex := MultiplexerChatCompletionRequest{
		ChatHistory: chatHistory,
		ModelName:   model_name.ModelName,
		Functions:   []FunctionCall{},
	}
	return createComplexRequest(requestMultiplex)
}

var rag_query_func_schema = openai.FunctionDefinition{
	Name: "query_government_documents",
	Parameters: jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"query": {
				Type:        jsonschema.String,
				Description: "The query string to search government documents and knowledge",
			},
		},
		Required: []string{"query"},
	},
}

// arguments='{"order_id":"order_12345"}',
func rag_query_func_generated_from_filters(filters search.Metadata) func(query_json string) (ToolCallResults, error) {
	return func(query_json string) (ToolCallResults, error) {
		var queryData map[string]string
		err := json.Unmarshal([]byte(query_json), &queryData)
		if err != nil {
			return ToolCallResults{}, fmt.Errorf("error unmarshaling query_json: %v", err)
		}
		search_query, ok := queryData["query"]
		if !ok {
			return ToolCallResults{}, fmt.Errorf("query field is missing in query_json")
		}
		search_request := search.SearchRequest{search_query, filters}
		search_results, err := search.SearchQuickwit(search_request)
		if err != nil {
			return ToolCallResults{}, err
		}
		// Increase to give llm more results.
		const truncation = 4
		var truncated_search_results []search.SearchData
		if len(search_results) < truncation {
			truncated_search_results = search_results
		} else {
			truncated_search_results = search_results[:truncation]
		}
		format_string := search.FormatSearchResults(truncated_search_results, search_query)
		result := ToolCallResults{Response: format_string, Citations: &truncated_search_results}

		return result, nil
	}
}

func rag_func_call_filters(filters search.Metadata) FunctionCall {
	return FunctionCall{
		Schema: rag_query_func_schema,
		Func:   rag_query_func_generated_from_filters(filters),
	}
}

var rag_func_call_no_filters = rag_func_call_filters(search.Metadata{})

func (model_name LLMModel) RagChat(chatHistory []ChatMessage, filters search.Metadata) (ChatMessage, error) {
	requestMultiplex := MultiplexerChatCompletionRequest{
		ChatHistory: chatHistory,
		ModelName:   model_name.ModelName,
		Functions:   []FunctionCall{rag_func_call_filters(filters)},
	}
	return createComplexRequest(requestMultiplex)
}

type LLM interface {
	Chat(chatHistory []ChatMessage) (ChatMessage, error)
	RagChat(chatHistory []ChatMessage) (ChatMessage, error)
}

// Wait to add all the llm utils until you understand how to write concurrent code in go more.
