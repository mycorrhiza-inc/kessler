package search

import (
	"context"
	"fmt"
	"kessler/objects/networking"
	"kessler/routing"
)

func SearchDataPassesFilters(result SearchDataHydrated, filters networking.FilterFields) bool {
	docket_id := filters.MetadataFilters.DocketID
	if result.File.Conversation.DocketID != docket_id {
		return false
	}
	return true
}

func ValidateSearchResultsAgainstFilters(ctx context.Context, results []SearchData, filters networking.FilterFields) error {
	q := *routing.DBQueriesFromContext(ctx)
	hydrated_results, err := HydrateSearchResults(results, ctx, q)
	if err != nil {
		return err
	}
	return ValidateSearchResultsHydratedAgainstFilters(hydrated_results, filters)
}

func ValidateSearchResultsHydratedAgainstFilters(results []SearchDataHydrated, filters networking.FilterFields) error {
	invalid_count := 0
	for i, result := range results {
		if !SearchDataPassesFilters(result, filters) {
			invalid_count += 1
			results[i].Name = "INVALID RESULT DUE TO FILTERS: " + results[i].Name
		}
	}
	if invalid_count > 0 {
		return fmt.Errorf("Found %v invalid results out of %v total results", invalid_count, len(results))
	}
	return nil
}
