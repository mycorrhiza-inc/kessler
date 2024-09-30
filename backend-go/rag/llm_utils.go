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
	Model       string        `json:"model"`
	ChatHistory []ChatMessage `json:"chat_history"`
}

// write two functions, one that converts a simple chatmessage into a ChatMessage with emtpy values for context and citations, while validating the chatRole. Then write another that just throws that info away and turns a ChatMessage into a simple chatmessage.
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

type ChatMessage struct {
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

func SimpleToChatMessage(msg SimpleChatMessage) (ChatMessage, error) {
	role, err := ValidateChatRole(msg.Role)
	if err != nil {
		return ChatMessage{}, err
	}
	return ChatMessage{
		Content:   msg.Content,
		Role:      role,
		Citations: &[]SearchData{},
		Context:   &[]SimpleChatMessage{},
	}, nil
}

func KeToSimpleChatMessage(keMsg ChatMessage) SimpleChatMessage {
	return SimpleChatMessage{
		Content: keMsg.Content,
		Role:    string(keMsg.Role),
	}
}

func SimpleToChatMessages(msgs []SimpleChatMessage) ([]ChatMessage, error) {
	var keMsgs []ChatMessage
	for _, msg := range msgs {
		keMsg, err := SimpleToChatMessage(msg)
		if err != nil {
			return nil, err
		}
		keMsgs = append(keMsgs, keMsg)
	}
	return keMsgs, nil
}

func KeToSimpleChatMessages(keMsgs []ChatMessage) []SimpleChatMessage {
	var msgs []SimpleChatMessage
	for _, keMsg := range keMsgs {
		msg := KeToSimpleChatMessage(keMsg)
		msgs = append(msgs, msg)
	}
	return msgs
}

func CreateKeChatCompletion(modelName string, chatHistory []ChatMessage) (ChatMessage, error) {
	simple_history := KeToSimpleChatMessages(chatHistory)
	simple_completion_string, err := createSimpleChatCompletionString(modelName, simple_history)
	if err != nil {
		return ChatMessage{}, err
	}
	ke_completion := ChatMessage{
		simple_completion_string,
		Assistant,
		&[]SearchData{},
		&[]SimpleChatMessage{},
	}
	return ke_completion, nil
}

type LLMModel struct {
	model_name string
}

func (model_name LLMModel) Achat(chatHistory []ChatMessage) (ChatMessage, error) {
	return CreateKeChatCompletion(model_name.model_name, chatHistory)
}

type LLM interface {
	Achat(chatHistory []ChatMessage) (ChatMessage, error)
}

func SummarizeSingleChunk(model LLM, markdownText string) (string, error) {
	const summarizePrompt = "Make sure to provide a well researched summary of the text provided by the user, if it appears to be the summary of a larger document, just summarize the section provided."
	summarizeMessage := ChatMessage{
		Role:    System,
		Content: summarizePrompt,
	}
	textMessage := ChatMessage{
		Role:    User,
		Content: markdownText,
	}
	history := []ChatMessage{summarizeMessage, textMessage}
	summary, err := model.Achat(history)
	if err != nil {
		return "", err
	}
	return summary.Content, nil
}

func SimpleInstruct(model LLM, content string, instruct string) (string, error) {
	history := []ChatMessage{
		{Content: instruct, Role: System},
		{Content: content, Role: User},
	}
	completion, err := model.Achat(history)
	if err != nil {
		return "", err
	}
	return completion.Content, nil
}

// func SummarizeMapReduce(model LLM, markdownText string, maxTokenSize int) (string, error) {
// 	splits := splitByMaxTokenSize(markdownText, maxTokenSize)
// 	if len(splits) == 1 {
// 		return SummarizeSingleChunk(model, markdownText)
// 	}
//
// 	var summaries []string
// 	for _, chunk := range splits {
// 		summary, err := SummarizeSingleChunk(model, chunk)
// 		if err != nil {
// 			return "", err
// 		}
// 		summaries = append(summaries, summary)
// 	}
//
// 	const coherencePrompt = "Please rewrite the following list of summaries of chunks of the document into a final summary of similar length that incorporates all the details present in the chunks"
// 	cohereMessage := ChatMessage{
// 		Role:    System,
// 		Content: coherencePrompt,
// 	}
// 	combinedSummariesPrompt := ChatMessage{
// 		Role:    User,
// 		Content: strings.Join(summaries, "\n"),
// 	}
// 	finalSummary, err := model.Achat([]ChatMessage{cohereMessage, combinedSummariesPrompt})
// 	if err != nil {
// 		return "", err
// 	}
// 	return finalSummary.Content, nil
// }
//
// func MapReduceLLMInstructionAcrossString(model LLM, content string, chunkSize int, instruction string, joinStr string) (string, error) {
// 	splits := tokenSplit(content, chunkSize)
//
// 	type result struct {
// 		content string
// 		err     error
// 	}
//
// 	resultChan := make(chan result, len(splits))
//
// 	for _, chunk := range splits {
// 		go func(chunk string) {
// 			history := []ChatMessage{
// 				{Content: instruction, Role: System},
// 				{Content: chunk, Role: User},
// 			}
// 			completion, err := model.Achat(history)
// 			resultChan <- result{content: completion.Content, err: err}
// 		}(chunk)
// 	}
//
// 	var results []string
// 	for i := 0; i < len(splits); i++ {
// 		res := <-resultChan
// 		if res.err != nil {
// 			return "", res.err
// 		}
// 		results = append(results, res.content)
// 	}
//
// 	return strings.Join(results, joinStr), nil
// }
