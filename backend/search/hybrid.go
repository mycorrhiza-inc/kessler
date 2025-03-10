package search

import (
	"fmt"

	"go.uber.org/zap"
)

type SearchReturn struct {
	Results []SearchDataHydrated
	Error   error
}

func HybridSearch(request SearchRequest) ([]SearchDataHydrated, error) {
	// This could absolutely be dryified
	chanMilvus := make(chan SearchReturn)
	chanQuickwit := make(chan SearchReturn)
	go func() {
		// Not implemented
		// results, err := SearchMilvus(request)
		err := fmt.Errorf("not implemented")
		return_results := SearchReturn{[]SearchDataHydrated{}, err}
		chanMilvus <- return_results
	}()
	go func() {
		results, err := SearchQuickwit(request)
		return_results := SearchReturn{results, err}
		chanQuickwit <- return_results
	}()
	resultsMilvus := <-chanMilvus
	resultsQuickwit := <-chanQuickwit
	if resultsMilvus.Error == nil {
		log.Info("Milvus returned error", zap.Error(resultsMilvus.Error))
	}
	if resultsQuickwit.Error == nil {
		log.Info("Quickwit returned error", zap.Error(resultsQuickwit.Error))
	}
	if resultsMilvus.Error != nil && resultsQuickwit.Error != nil {
		return []SearchDataHydrated{}, fmt.Errorf("both Milvus and Quickwit returned errors. milvus error: %s quickwit error: %s", resultsMilvus.Error, resultsQuickwit.Error)
	}
	unrankedResults := append(resultsMilvus.Results, resultsQuickwit.Results...)
	rerankedData, err := rerankSearchResults(unrankedResults, request.Query)
	// Fail semi silently and returns the regular unranked results
	if err != nil {
		log.Info("Error reranking results", zap.Error(err))
		return unrankedResults, nil
	}

	return rerankedData, nil
}
