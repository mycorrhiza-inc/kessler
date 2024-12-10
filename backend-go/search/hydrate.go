package search

import (
	"context"
	"kessler/objects/files"

	"fmt/log"
)

func HydrateSearchResults(results []SearchData) []SearchData {
	// get a pg connection
	// query the database for all metadata of each result
	// return array of hydrated results
	// This could be done in a more efficient way
	return results
}

func SearchThenHydrate(r SearchRequest, ctx context.Context) ([]files.FileMetadataSchema, error) {

	data, err := SearchQuickwit(r)
	if err != nil {
		log.Errorf("search-then-hydrate: error getting quickwit data", err)
	}
	idList := make([]string, len(data))
	for i, d := range data {
		idList[i] = d.SourceID
	}

	hydratedData := []files.FileMetadataSchema{}

	return []files.FileMetadataSchema{}, nil
}
