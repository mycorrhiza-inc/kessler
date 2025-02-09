package search

import (
	"fmt"
	"kessler/objects/networking"
	"testing"

	"github.com/charmbracelet/log"
)

var example_search_request = SearchRequest{
	Index:         "Marshall Fire",
	Query:         "",
	SearchFilters: networking.FilterFields{},
	SortBy:        []string{},
	MaxHits:       10,
	StartOffset:   0,
	GetText:       false,
}

func TestQuickwit(t *testing.T) {
	results, err := SearchQuickwit(example_search_request)
	if err != nil {
		log.Info("Error:", err)
		return
	}
	resultsString := fmt.Sprintf("%v", results)
	log.Info("Results:", resultsString)
}
