package search

import (
	"fmt"
	"testing"
)

var example_search_request = SearchRequest{
	Index:         "Marshall Fire",
	Query:         "",
	SearchFilters: FilterFields{},
	SortBy:        []string{},
	MaxHits:       10,
	StartOffset:   0,
	GetText:       false,
}

func TestQuickwit(t *testing.T) {
	results, err := SearchQuickwit(example_search_request)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	resultsString := fmt.Sprintf("%v", results)
	fmt.Println("Results:", resultsString)
}
