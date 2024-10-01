package rag

import (
	"encoding/json"
	"fmt"

	"github.com/mycorrhiza-inc/kessler/backend-go/search"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

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
func rag_query_func(query_json string) (string, error) {
	var queryData map[string]string
	err := json.Unmarshal([]byte(query_json), &queryData)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling query_json: %v", err)
	}
	search_query, ok := queryData["query"]
	if !ok {
		return "", fmt.Errorf("query field is missing in query_json")
	}
	search_request := search.SearchRequest{search_query, search.Metadata{}}
	search_results, err := search.SearchQuickwit(search_request)
	if err != nil {
		return "", err
	}
	format_string := search.FormatSearchResults(search_results, 1)

	return format_string, nil
}
