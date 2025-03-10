package rag

import (
	"context"
	"encoding/json"
	"fmt"
	"kessler/internal/llm_utils"
	"kessler/internal/objects/networking"
	"kessler/internal/search"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type RagLLMModel llm_utils.LLMModel

func (ragLlm RagLLMModel) RagChat(ctx context.Context, chatHistory []llm_utils.ChatMessage, filters networking.FilterFields) (llm_utils.ChatMessage, error) {
	// TODO: Add a method to the LLMModel interface to get the chat history from the OpenAI API
	regular_llm := llm_utils.LLMModel(ragLlm)
	result, err := regular_llm.Chat(ctx, chatHistory)
	return result, err
}

var rag_query_func_schema = openai.FunctionDefinition{
	Name:        "query_government_documents",
	Description: "Search the text of goverment documents to find snippets from documents that match.",
	Parameters: jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"query": {
				Type:        jsonschema.String,
				Description: "The query string to search.",
			},
		},
		Required: []string{"query"},
	},
}

func rag_func_call_filters(filters networking.FilterFields) llm_utils.FunctionCall {
	return llm_utils.FunctionCall{
		Schema: rag_query_func_schema,
		Func: func(query_json string) (llm_utils.ToolCallResults, error) {
			var queryData map[string]string
			err := json.Unmarshal([]byte(query_json), &queryData)
			if err != nil {
				return llm_utils.ToolCallResults{}, fmt.Errorf("error unmarshaling query_json: %v", err)
			}
			search_query, ok := queryData["query"]
			if !ok {
				return llm_utils.ToolCallResults{}, fmt.Errorf("query field is missing in query_json")
			}
			return rag_query_func_generated_from_filters(filters)(search_query)
		},
	}
}

func rag_query_func_generated_from_filters(filters networking.FilterFields) func(search_query string) (llm_utils.ToolCallResults, error) {
	return func(search_query string) (llm_utils.ToolCallResults, error) {
		search_request := search.SearchRequest{
			Index:         "NY_PUC",
			Query:         search_query,
			SearchFilters: filters,
			GetText:       true,
			MaxHits:       30,
		}
		search_results, err := search.HybridSearch(search_request)
		if err != nil {
			return llm_utils.ToolCallResults{}, err
		}
		// Increase to give llm more results.
		const truncation = 4
		var truncated_search_results []search.SearchDataHydrated
		if len(search_results) < truncation {
			truncated_search_results = search_results
		} else {
			truncated_search_results = search_results[:truncation]
		}
		format_string := search.FormatSearchResults(truncated_search_results, search_query)
		result := llm_utils.ToolCallResults{Response: format_string}
		// result := llm_utils.ToolCallResults{Response: format_string, Citations: []llm_utils.Citation(truncated_search_results)}

		return result, nil
	}
}
