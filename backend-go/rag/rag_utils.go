package rag

import (
	"fmt"
	"strings"
)

func DoesChatNeedQuery(model LLM, chatHistory []KeChatMessage) (bool, error) {
	const doesChatNeedQuery = "Please determine if you need to query a vector database of relevant documents to answer the user. Answer with only a \"yes\" or \"no\"."
	checkMessage := KeChatMessage{
		Role:    Assistant,
		Content: doesChatNeedQuery,
	}
	checkHistory := append(chatHistory, checkMessage)
	checkResponse, err := model.Achat(checkHistory)
	if err != nil {
		return true, err
	}

	checkYesNo := func(testStr string) (bool, error) {
		testStr = strings.ToLower(testStr)
		if strings.HasPrefix(testStr, "yes") {
			return true, nil
		}
		if strings.HasPrefix(testStr, "no") {
			return false, nil
		}
		return true, fmt.Errorf("expected yes or no, got: %s", testStr)
	}

	return checkYesNo(checkResponse.Content)
}

// func RagAchat(model LLM, chatHistory []KeChatMessage, filesRepo FileRepository, logger *log.Logger) (KeChatMessage, []FileSchema, error) {
// 	if logger == nil {
// 		logger = log.Default()
// 	}
// 	needsQuery, err := model.DoesChatNeedQuery(chatHistory)
// 	if err != nil {
// 		return KeChatMessage{}, nil, err
// 	}
// 	if !needsQuery {
// 		finalMessage, err := model.Achat(chatHistory)
// 		return finalMessage, nil, err
// 	}
//
// 	generateQueryFromChatHistory := func(chatHistory []KeChatMessage) (string, error) {
// 		querygenAddendum := KeChatMessage{
// 			Role:    System,
// 			Content: generateQueryFromChatHistoryPrompt,
// 		}
// 		completion, err := model.Achat(append(chatHistory, querygenAddendum))
// 		if err != nil {
// 			return "", err
// 		}
// 		return completion.Content, nil
// 	}
//
// 	generateContextMsgFromSearchResults := func(searchResults []map[string]interface{}, maxResults int) KeChatMessage {
// 		if logger == nil {
// 			logger = log.Default()
// 		}
// 		res := searchResults
// 		if len(res) > maxResults {
// 			res = res[:maxResults]
// 		}
// 		const returnPrompt := "Here is a list of documents that might be relevant to the following chat:"
// 		for _, result := range res {
// 			uuidStr := result["entity"].(map[string]interface{})["source_id"].(string)
// 			text := result["entity"].(map[string]interface{})["text"].(string)
// 			returnPrompt += fmt.Sprintf("\n\n%s:\n%s", uuidStr, text)
// 		}
// 		return KeChatMessage{
// 			Role:    ChatRoleAssistant,
// 			Content: returnPrompt,
// 		}
// 	}
//
// 	query, err := generateQueryFromChatHistory(chatHistory)
// 	if err != nil {
// 		return KeChatMessage{}, nil, err
// 	}
// 	res := search(query, []string{"source_id", "text"})
// 	logger.Println(res)
// 	contextMsg := generateContextMsgFromSearchResults(res, 3)
//
// 	finalMessage, err := model.Achat(append([]KeChatMessage{contextMsg}, chatHistory...))
// 	if err != nil {
// 		return KeChatMessage{}, nil, err
// 	}
// 	returnSchemas, err := convertSearchResultsToFrontendTable(res, filesRepo)
// 	if err != nil {
// 		return KeChatMessage{}, nil, err
// 	}
// 	finalMessage.Citations = returnSchemas
//
// 	return finalMessage, returnSchemas, nil
// }
