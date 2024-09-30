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
type LLM interface {
	Achat(chatHistory []KeChatMessage) (KeChatMessage, error)
}
		Role:    string(keMsg.Role),
	}
}

func SimpleToKeChatMessages(msgs []SimpleChatMessage) ([]KeChatMessage, error) {
	var keMsgs []KeChatMessage
	for _, msg := range msgs {
		keMsg, err := SimpleToKeChatMessage(msg)
		if err != nil {
			return nil, err
		}
		keMsgs = append(keMsgs, keMsg)
	}
	return keMsgs, nil
}

func KeToSimpleChatMessages(keMsgs []KeChatMessage) []SimpleChatMessage {
	var msgs []SimpleChatMessage
	for _, keMsg := range keMsgs {
		msg := KeToSimpleChatMessage(keMsg)
		msgs = append(msgs, msg)
	}
	return msgs
}

func CreateKeChatCompletion(modelName string, chatHistory []KeChatMessage) (KeChatMessage, error) {
	simple_history := KeToSimpleChatMessages(chatHistory)
	simple_completion_string, err := createSimpleChatCompletionString(modelName, simple_history)
	if err != nil {
		return KeChatMessage{}, err
	}
	ke_completion := KeChatMessage{
		simple_completion_string,
		Assistant,
		&[]SearchData{},
		&[]SimpleChatMessage{},
	}
	return ke_completion, nil
}
type LLMModelName struct {
	model_name string
}

func (model_name LLMModelName) Achat(chatHistory []KeChatMessage) (KeChatMessage, error) {
	return CreateKeChatCompletion(model_name, chatHistory)
}
type LLM interface{
	Achat(chatHistory []KeChatMessage) (KeChatMessage, error)
}



func (model LLM) SummarizeSingleChunk(markdownText string) (string, error) {
	summarizePrompt := "Make sure to provide a well researched summary of the text provided by the user, if it appears to be the summary of a larger document, just summarize the section provided."
	summarizeMessage := KeChatMessage{
		Role:    ChatRoleSystem,
		Content: summarizePrompt,
	}
	textMessage := KeChatMessage{
		Role:    ChatRoleUser,
		Content: markdownText,
	}
	history := []KeChatMessage{summarizeMessage, textMessage}
	summary, err := model.Achat(history)
	if err != nil {
		return "", err
	}
	return summary.Content, nil
}

func (model LLM) SummarizeMapReduce(markdownText string, maxTokenSize int) (string, error) {
	splits := splitByMaxTokenSize(markdownText, maxTokenSize)
	if len(splits) == 1 {
		return model.SummarizeSingleChunk(markdownText)
	}

	var summaries []string
	for _, chunk := range splits {
		summary, err := model.SummarizeSingleChunk(chunk)
		if err != nil {
			return "", err
		}
		summaries = append(summaries, summary)
	}

	coherencePrompt := "Please rewrite the following list of summaries of chunks of the document into a final summary of similar length that incorporates all the details present in the chunks"
	cohereMessage := KeChatMessage{
		Role:    ChatRoleSystem,
		Content: coherencePrompt,
	}
	combinedSummariesPrompt := KeChatMessage{
		Role:    ChatRoleUser,
		Content: strings.Join(summaries, "\n"),
	}
	finalSummary, err := model.Achat([]KeChatMessage{cohereMessage, combinedSummariesPrompt})
	if err != nil {
		return "", err
	}
	return finalSummary.Content, nil
}

func (model LLM) SimpleInstruct(content string, instruct string) (string, error) {
	history := []KeChatMessage{
		{Content: instruct, Role: ChatRoleSystem},
		{Content: content, Role: ChatRoleUser},
	}
	completion, err := model.Achat(history)
	if err != nil {
		return "", err
	}
	return completion.Content, nil
}

func (model LLM) MapReduceLLMInstructionAcrossString(content string, chunkSize int, instruction string, joinStr string) (string, error) {
	splits := tokenSplit(content, chunkSize)

	type result struct {
		content string
		err     error
	}

	resultChan := make(chan result, len(splits))

	for _, chunk := range splits {
		go func(chunk string) {
			history := []KeChatMessage{
				{Content: instruction, Role: ChatRoleSystem},
				{Content: chunk, Role: ChatRoleUser},
			}
			completion, err := model.Achat(history)
			resultChan <- result{content: completion.Content, err: err}
		}(chunk)
	}

	var results []string
	for i := 0; i < len(splits); i++ {
		res := <-resultChan
		if res.err != nil {
			return "", res.err
		}
		results = append(results, res.content)
	}

	return strings.Join(results, joinStr), nil
}
