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

type LLM interface {
	Achat(chatHistory []map[string]interface{}) (string, error)
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

// async def summarize_single_chunk(self, markdown_text: str) -> str:
//     summarize_prompt = "Make sure to provide a well researched summary of the text provided by the user, if it appears to be the summary of a larger document, just summarize the section provided."
//     summarize_message = KeChatMessage(
//         role=ChatRole.system, content=summarize_prompt
//     )
//     text_message = KeChatMessage(role=ChatRole.user, content=markdown_text)
//     summary = await self.achat(
//         sanitzie_chathistory_llamaindex([summarize_message, text_message])
//     )
//     return summary.content
//
// async def summarize_mapreduce(
//     self, markdown_text: str, max_tokensize: int = 8096
// ) -> str:
//     splits = split_by_max_tokensize(markdown_text, max_tokensize)
//     if len(splits) == 1:
//         return await self.summarize_single_chunk(markdown_text)
//     summaries = await asyncio.gather(
//         *[self.summarize_single_chunk(chunk) for chunk in splits]
//     )
//     coherence_prompt = "Please rewrite the following list of summaries of chunks of the document into a final summary of similar length that incorperates all the details present in the chunks"
//     cohere_message = KeChatMessage(role=ChatRole.system, content=coherence_prompt)
//     combined_summaries_prompt = KeChatMessage(
//         role=ChatRole.user, content="\n".join(summaries)
//     )
//     final_summary = await self.achat([cohere_message, combined_summaries_prompt])
//     return final_summary.content
//
// async def simple_instruct(self, content: str, instruct: str) -> str:
//     history = [
//         KeChatMessage(content=instruct, role=ChatRole.system),
//         KeChatMessage(content=content, role=ChatRole.user),
//     ]
//     completion = await self.achat(history)
//     return completion.content
//
// async def mapreduce_llm_instruction_across_string(
//     self, content: str, chunk_size: int, instruction: str, join_str: str
// ) -> str:
//     # Replace with semantic splitter
//     split = token_split(content, chunk_size)
//
//     async def clean_chunk(chunk: str) -> str:
//
//         history = [
//             KeChatMessage(content=instruction, role=ChatRole.system),
//             KeChatMessage(content=chunk, role=ChatRole.user),
//         ]
//         completion = await self.llm.achat(history)
//         return completion.content
//
//     tasks = [clean_chunk(chunk) for chunk in split]
//     results = await asyncio.gather(*tasks)
//     return join_str.join(results)
//
