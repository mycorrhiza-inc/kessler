package search

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type SearchRequest struct {
	Query          string   `json:"query"`
	SearchFilters Metadata `json:"filters"`
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
