package search

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type SearchRequest struct {
	Query string `json:"query"`
}

func HandleSearchRequest(w http.ResponseWriter, r *http.Request) {
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

	data, err := searchQuickwit(string(RequestData.Query))
	if err != nil {
		log.Fatalf("Error searching quickwit: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respString, err := json.Marshal(data)

	if err != nil {
		log.Fatal("Error marshalling response data")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(respString))
}
