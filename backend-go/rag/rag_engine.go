package rag

import (
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
	search_query := "get query from json when query_json is formatted like {\"query\": \"what is the capital of the united states\"}"
	search_request := search.SearchRequest{search_query, search.Metadata{}}
	search_results, err := search.SearchQuickwit(search_request)
	if err != nil {
		return "", err
	}
	format_string := search.FormatSearchResults(search_results, 1)

	return format_string, nil
}
