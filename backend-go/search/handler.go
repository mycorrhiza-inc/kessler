package search

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type SearchRequest struct {
	Index         string       `json:"index"`
	Query         string       `json:"query"`
	SearchFilters FilterFields `json:"filters"`
	SortBy        []string     `json:"sort_by"`
	MaxHits       int          `json:"max_hits"`
	StartOffset   int          `json:"start_offset"`
	GetText       bool         `json:"get_text"`
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
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}
		defer r.Body.Close() // Close the body when done

		data, err := SearchQuickwit(RequestData)
		if err != nil {
			log.Printf("Error searching quickwit: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		respString, err := json.Marshal(data)

		if err != nil {
			log.Println("Error marshalling response data")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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
		data, err := GetRecentCaseData(0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respString, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(respString))
	case http.MethodPost:
		var request RecentUpdatesRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}
		data, err := GetRecentCaseData(request.Page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respString, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(respString))
	default:
		http.Error(w, "Unsupported request method", http.StatusMethodNotAllowed)

	}
}
