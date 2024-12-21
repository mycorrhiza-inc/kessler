package search

import (
	"context"
	"fmt"
	"kessler/objects/networking"
	"kessler/routing"
	"log"
)

func SearchDataPassesFilters(result SearchDataHydrated, filters networking.FilterFields) bool {
	// mdata := result.File.Mdata
	docket_id := filters.MetadataFilters.DocketID
	if docket_id != "" {
		result_docket_id := result.File.Mdata["docket_id"]
		if result_docket_id != docket_id {
			log.Printf("Docket ID mismatch, wanted: %v \ngot: %v\n", docket_id, result_docket_id)
			return false
		}
	}
	file_class := filters.MetadataFilters.FileClass
	if file_class != "" {
		result_file_class := result.File.Mdata["file_class"]
		if result_file_class != file_class {
			log.Printf("File Class mismatch, wanted: %v \ngot: %v\n", file_class, result_file_class)
			return false
		}
	}

	// doctype := filters.MetadataFilters.Doctype
	// if result.File.Extension != doctype {
	// 	return false
	// }
	return true
}

func ValidateAgainstFilters(ctx context.Context, results []SearchData, filters networking.FilterFields) ([]SearchDataHydrated, error) {
	q := *routing.DBQueriesFromContext(ctx)
	hydrated_results, err := HydrateSearchResults(results, ctx, q)
	if err != nil {
		return []SearchDataHydrated{}, err
	}
	return ValidateHydratedAgainstFilters(hydrated_results, filters)
}

func ValidateHydratedAgainstFilters(results []SearchDataHydrated, filters networking.FilterFields) ([]SearchDataHydrated, error) {
	invalid_count := 0
	for i, result := range results {
		if !SearchDataPassesFilters(result, filters) {
			invalid_count += 1
			results[i].Name = "INVALID RESULT DUE TO FILTERS: " + results[i].Name
		}
	}
	if invalid_count > 0 {
		log.Printf("Filters invalid results: %v out of %v total results", invalid_count, len(results))
		return results, fmt.Errorf("Found %v invalid results out of %v total results", invalid_count, len(results))
	}
	return results, nil
}
