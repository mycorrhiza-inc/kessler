package search

import (
	"context"
	"kessler/objects/files"

	"log"
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
		log.Printf("search-then-hydrate: error getting quickwit data", err)
	}
	idList := make([]string, len(data))
	for i, d := range data {
		idList[i] = d.SourceID
	}

	// hydratedData := make([]files.FileMetadataSchema, len(idList))

	return []files.FileMetadataSchema{}, nil
}
