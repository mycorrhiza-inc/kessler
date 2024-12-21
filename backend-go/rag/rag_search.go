package rag

import (
	"encoding/json"
	"fmt"
	"kessler/objects/networking"
	"kessler/search"
)

func rag_query_func_generated_from_filters(filters networking.FilterFields) func(query_json string) (ToolCallResults, error) {
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
		search_request := search.SearchRequest{
			Index:         "NY_PUC",
			Query:         search_query,
			SearchFilters: filters,
			GetText:       true,
			MaxHits:       30,
		}
		search_results, err := search.HybridSearch(search_request)
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
