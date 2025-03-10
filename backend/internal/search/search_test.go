package search

import (
	"fmt"
	"kessler/internal/objects/networking"
	"testing"

	"go.uber.org/zap"
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
		log.Info("search error", zap.Error(err))
		return
	}
	resultsString := fmt.Sprintf("%v", results)
	log.Info("search results", zap.String("results", resultsString))
}
