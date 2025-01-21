package search

import (
	"encoding/json"
	"fmt"
	"kessler/objects/networking"
	"log"
	"net/http"
	"strings"
)

type SearchRequest struct {
	Index         string                  `json:"index"`
	Query         string                  `json:"query"`
	SearchFilters networking.FilterFields `json:"filters"`
	SortBy        []string                `json:"sort_by"`
	MaxHits       int                     `json:"max_hits"`
	StartOffset   int                     `json:"start_offset"`
	GetText       bool                    `json:"get_text"`
}

func (s SearchRequest) String() string {
	return fmt.Sprintf(
		"SearchRequest: {\n\tIndex: %s\n\tQuery: %s\n\tFilters: %s\n\tSortBy: %s\n\tMaxHits: %d\n\tStartOffset: %d\nGetText: %t\n}\n",
		s.Index,
		s.Query,
		s.SearchFilters,
		strings.Join(s.SortBy, ", "),
		s.MaxHits,
		s.StartOffset,
		s.GetText,
	)
}

type RecentUpdatesRequest struct {
	Page int `json:"page"`
}

func HandleSearchRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprint(w, "Hi there!")
		return
	case http.MethodPost:
		log.Println("Received a search request")

		// Create an instance of RequestData
		var RequestData SearchRequest

		// Decode the JSON body into the struct
		err := json.NewDecoder(r.Body).Decode(&RequestData)
		if err != nil {
			log.Printf("Error decoding JSON: %v\n", err)
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}
		defer r.Body.Close() // Close the body when done

		pagination := networking.PaginationFromUrlParams(r)
		RequestData.MaxHits = int(pagination.Limit)
		RequestData.StartOffset = int(pagination.Offset)

		hydrated_data, err := SearchQuickwit(RequestData)
		if err != nil {
			log.Printf("Error searching quickwit: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// q := *routing.DBQueriesFromRequest(r)
		// ctx := r.Context()
		// hydrated_data, err := HydrateSearchResults(data, ctx, q)
		// if err != nil {
		// 	errorstring := fmt.Sprintf("Error hydrating search results: %v\n", err)
		// 	log.Println(errorstring)
		// 	http.Error(w, errorstring, http.StatusInternalServerError)
		// 	return
		// }
		// TODO : Reneable validation once other stuff is certainly working.
		_, err = ValidateHydratedAgainstFilters(hydrated_data, RequestData.SearchFilters)
		if err != nil {
			errorstring := fmt.Sprintf("Returned results did not match filters: %v\n", err)
			log.Println(errorstring)
			http.Error(w, errorstring, http.StatusInternalServerError)
			return
		}

		respString, err := json.Marshal(hydrated_data)
		if err != nil {
			log.Println("Error marshalling response data")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(respString))

	case http.MethodPut:
		fmt.Fprintf(w, "PUT request")
	case http.MethodDelete:
		fmt.Fprintf(w, "DELETE request")
	default:
		http.Error(w, "Unsupported request method", http.StatusMethodNotAllowed)
	}
}

func HandleRecentUpdatesRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		pagination := networking.PaginationFromUrlParams(r)
		maxHits := int(pagination.Limit)
		offset := int(pagination.Offset)

		hydrated_data, err := GetRecentCaseData(maxHits, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// q := *routing.DBQueriesFromRequest(r)
		// ctx := r.Context()
		// hydrated_data, err := HydrateSearchResults(data, ctx, q)
		// if err != nil {
		// 	errorstring := fmt.Sprintf("Error hydrating search results: %v\n", err)
		// 	log.Println(errorstring)
		// 	http.Error(w, errorstring, http.StatusInternalServerError)
		// 	return
		// }
		respString, err := json.Marshal(hydrated_data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(respString))
	default:
		http.Error(w, "Unsupported request method", http.StatusMethodNotAllowed)

	}
}
