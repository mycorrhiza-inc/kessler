package search

import (
	"context"
	"fmt"
	"kessler/objects/networking"
	"kessler/routing"
	"log"
	"reflect"
	"strings"
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

// write a function that will take in a searchRequest and searchResults and create a new quickwitSearchResponse then for each hit and snippet in the passed in search results, make sure all the filters in search request are valid for that, if it is valid append it to the return searchResponse, else skip it and print a scary error message, then return the list of validated results
func ValidateSearchRequest(searchRequest SearchRequest, searchResults quickwitSearchResponse) quickwitSearchResponse {
	global_filters := searchRequest.SearchFilters
	filters := global_filters.MetadataFilters
	metadata_filters := filters.Metadata
	var validatedResponse quickwitSearchResponse

	for i, hit := range searchResults.Hits {
		isValid := true

		// Validate query matches if present
		if searchRequest.Query != "" {
			if !strings.Contains(strings.ToLower(hit.Text), strings.ToLower(searchRequest.Query)) &&
				!strings.Contains(strings.ToLower(hit.Name), strings.ToLower(searchRequest.Query)) {
				log.Printf("❌ WARNING: Hit %d failed query validation", i)
				isValid = false
			}
		}

		// Claude writes some interesting go code lol - nic
		// Validate metadata filters
		v := reflect.ValueOf(metadata_filters)
		t := reflect.TypeOf(metadata_filters)
		for j := 0; j < t.NumField(); j++ {
			field := t.Field(j)
			value := v.Field(j)

			// Skip empty string fields
			if value.Kind() == reflect.String && value.String() != "" {
				hitValue := reflect.ValueOf(hit.Metadata).FieldByName(field.Name)
				if !hitValue.IsValid() || hitValue.String() != value.String() {
					log.Printf("❌ WARNING: Hit %d failed metadata validation for field %s", i, field.Name)
					isValid = false
					break
				}
			}
		}

		// If all validations pass, append to response
		if isValid {
			validatedResponse.Hits = append(validatedResponse.Hits, hit)
			if i < len(searchResults.Snippets) {
				validatedResponse.Snippets = append(validatedResponse.Snippets, searchResults.Snippets[i])
			}
		}
	}

	return validatedResponse
}
