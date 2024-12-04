package search

import (
	"fmt"
	"log"
)

type SearchReturn struct {
	Results []SearchData
	Error   error
}

func HybridSearch(request SearchRequest) ([]SearchData, error) {
	// This could absolutely be dryified
	chanMilvus := make(chan SearchReturn)
	chanQuickwit := make(chan SearchReturn)
	go func() {
		results, err := SearchMilvus(request)
		return_results := SearchReturn{results, err}
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
		fmt.Printf("Milvus returned error: %s", resultsMilvus.Error)
	}
	if resultsQuickwit.Error == nil {
		fmt.Printf("Quickwit returned error: %s", resultsQuickwit.Error)
	}
	if resultsMilvus.Error != nil && resultsQuickwit.Error != nil {
		return []SearchData{}, fmt.Errorf("both Milvus and Quickwit returned errors. milvus error: %s quickwit error: %s", resultsMilvus.Error, resultsQuickwit.Error)
	}
	unrankedResults := append(resultsMilvus.Results, resultsQuickwit.Results...)
	rerankedData, err := rerankSearchResults(unrankedResults, request.Query)
	// Fail semi silently and returns the regular unranked results
	if err != nil {
		log.Printf("Error reranking results: %s", err)
		return unrankedResults, nil
	}

	return rerankedData, nil
}
