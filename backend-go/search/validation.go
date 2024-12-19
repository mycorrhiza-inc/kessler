package search

import (
	"context"
	"fmt"
	"kessler/objects/networking"
	"kessler/routing"
)

func SearchDataPassesFilters(result SearchDataHydrated, filters networking.FilterFields) bool {
	// mdata := result.File.Mdata
	docket_id := filters.MetadataFilters.DocketID
	if result.File.Conversation.DocketID != docket_id {
		return false
	}
	// file_class := filters.MetadataFilters.FileClass
	// if mdata["file_class"] != file_class {
	// 	return false
	// }
	// doctype := filters.MetadataFilters.Doctype
	// if result.File.Extension != doctype {
	// 	return false
	// }

	return true
}

func ValidateSearchResultsAgainstFilters(ctx context.Context, results []SearchData, filters networking.FilterFields) ([]SearchDataHydrated, error) {
	q := *routing.DBQueriesFromContext(ctx)
	hydrated_results, err := HydrateSearchResults(results, ctx, q)
	if err != nil {
		return []SearchDataHydrated{}, err
	}
	return ValidateSearchResultsHydratedAgainstFilters(hydrated_results, filters)
}

func ValidateSearchResultsHydratedAgainstFilters(results []SearchDataHydrated, filters networking.FilterFields) ([]SearchDataHydrated, error) {
	invalid_count := 0
	for i, result := range results {
		if !SearchDataPassesFilters(result, filters) {
			invalid_count += 1
			results[i].Name = "INVALID RESULT DUE TO FILTERS: " + results[i].Name
		}
	}
	if invalid_count > 0 {
		return results, fmt.Errorf("Found %v invalid results out of %v total results", invalid_count, len(results))
	}
	return results, nil
}
